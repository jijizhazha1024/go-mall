package logic

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	order2 "jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/order/internal/mq/delay"
	"jijizhazha1024/go-mall/services/users/users"
	"time"

	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type orderCreateDTO struct {
	OrderID       string
	PreOrderID    string
	UserID        int64
	PaymentMethod order.PaymentMethod
	Address       *order2.OrderAddresses
	Items         []*checkout.CheckoutItem
	OrderItems    []*order2.OrderItems
	Amounts       *coupons.CalculateCouponResp
}
type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateOrder 创建订单
func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderRequest) (*order.OrderDetailResponse, error) {
	// --------------- 参数校验 ---------------
	if err := l.validateRequest(in); err != nil {
		return nil, err
	}
	// --------------- 绑定数据 ---------------
	dto, err := l.collectOrderData(in)
	if err != nil {
		l.Logger.Errorw("collect order data failed")
		return nil, err
	}

	orderValue := dto.ToOrderModel()
	orderValue.CouponId = in.CouponId
	res := &order.OrderDetailResponse{}
	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		orderSession := l.svcCtx.OrderModel.WithSession(session)
		// --------------- check ---------------
		orderRes, err := orderSession.CheckOrderExistByPreOrderId(ctx, dto.PreOrderID, int32(dto.UserID))
		if err != nil {
			l.Logger.Errorw("check order exist failed", append(l.logContext(dto), logx.Field("err", err))...)
			return err
		}
		if orderRes {
			res.StatusCode = code.OrderExist
			res.StatusMsg = code.OrderExistMsg
			return nil
		}
		// --------------- insert ---------------
		if _, err = orderSession.Insert(l.ctx, orderValue); err != nil {
			l.Logger.Errorw("insert order failed", append(l.logContext(dto), logx.Field("err", err))...)
			return err
		}
		// 插入订单项关联商品
		if err := l.svcCtx.OrderItemModel.BulkInsert(session, dto.OrderItems); err != nil {
			l.Logger.Errorw("insert order items failed", append(l.logContext(dto), logx.Field("err", err))...)
			return err
		}
		// 插入订单地址
		dto.Address.OrderId = dto.OrderID
		if _, err := l.svcCtx.OrderAddress.Insert(l.ctx, dto.Address); err != nil {
			l.Logger.Errorw("insert order address failed", append(l.logContext(dto), logx.Field("err", err))...)
			return err
		}
		return nil
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if res.StatusCode != code.Success {
		// 跳过，当前事务处理，确保幂等
		l.Logger.Infow("transaction aborted", l.logContext(dto)...)
		return res, nil
	}
	// --------------- 订单超时 ---------------
	if err := l.svcCtx.OrderDelayMQ.Product(&delay.OrderReq{
		OrderId: dto.OrderID,
		UserID:  int32(in.UserId),
	}); err != nil {
		l.Logger.Errorw("publish order delay message failed", l.logContext(dto)...)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return res, nil
}
func (l *CreateOrderLogic) validateRequest(in *order.CreateOrderRequest) error {
	if in.PreOrderId == "" || in.UserId == 0 || in.AddressId == 0 || in.PaymentMethod == 0 {
		return status.Error(codes.Aborted, "参数不合法")
	}
	return nil
}
func (l *CreateOrderLogic) collectOrderData(in *order.CreateOrderRequest) (*orderCreateDTO, error) {

	g, ctx := errgroup.WithContext(l.ctx)
	var dto = &orderCreateDTO{
		PreOrderID:    in.PreOrderId,
		UserID:        int64(in.UserId),
		PaymentMethod: in.PaymentMethod,
		OrderID:       l.generateOrderID(),
	}
	g.Go(func() error {
		// 获取订单详情
		checkoutDetail, err := l.svcCtx.CheckoutRpc.GetCheckoutDetail(ctx, &checkout.CheckoutDetailReq{
			PreOrderId: in.PreOrderId,
			UserId:     int32(in.UserId),
		})
		if err != nil {
			logx.Errorw("call rpc GetCheckoutDetail failed", append(l.logContext(dto), logx.Field("err", err))...)
			return status.Error(codes.Aborted, err.Error())
		}
		if checkoutDetail.StatusCode != code.Success {
			return status.Error(codes.Aborted, checkoutDetail.StatusMsg)
		}
		dto.OrderItems = convertToOrderItems(dto.OrderID, checkoutDetail.Data.Items)

		if in.CouponId == "" {
			dto.Amounts = &coupons.CalculateCouponResp{
				OriginAmount: checkoutDetail.Data.OriginalAmount,
				FinalAmount:  checkoutDetail.Data.FinalAmount,
			}
			return nil
		}
		// 计算优惠价格
		couponResp, err := l.svcCtx.CouponRpc.CalculateCoupon(ctx, &coupons.CalculateCouponReq{
			UserId:   int32(in.UserId),
			CouponId: in.CouponId,
			Items:    convertToCouponItems(checkoutDetail.Data.Items),
		})
		if err != nil {
			logx.Errorw("call rpc CalculateCoupon failed", append(l.logContext(dto), logx.Field("err", err))...)
			return status.Error(codes.Aborted, err.Error())
		}
		if couponResp.StatusCode != code.Success {
			return status.Error(codes.Aborted, couponResp.StatusMsg)
		}
		dto.Amounts = &coupons.CalculateCouponResp{
			OriginAmount:   couponResp.OriginAmount,
			DiscountAmount: couponResp.DiscountAmount,
			FinalAmount:    couponResp.FinalAmount,
		}
		return nil
	})
	g.Go(func() error {
		addressResp, err := l.svcCtx.UserRpc.GetAddress(ctx, &users.GetAddressRequest{
			UserId:    in.UserId,
			AddressId: in.AddressId,
		})
		if err != nil {
			l.Logger.Errorw("call rpc GetAddress failed", append(l.logContext(dto), logx.Field("err", err))...)
			return status.Error(codes.Aborted, err.Error())
		}
		if addressResp.StatusCode != code.Success {
			return status.Error(codes.Aborted, addressResp.StatusMsg)
		}
		dto.Address = &order2.OrderAddresses{
			AddressId:       uint64(in.AddressId),
			RecipientName:   addressResp.Data.RecipientName,
			PhoneNumber:     sql.NullString{String: addressResp.Data.PhoneNumber, Valid: addressResp.Data.PhoneNumber != ""},
			Province:        sql.NullString{String: addressResp.Data.Province, Valid: addressResp.Data.Province != ""},
			City:            addressResp.Data.City,
			DetailedAddress: addressResp.Data.DetailedAddress,
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return dto, nil

}
func (d *orderCreateDTO) ToOrderModel() *order2.Orders {
	return &order2.Orders{
		OrderId:        d.OrderID,
		PreOrderId:     d.PreOrderID,
		UserId:         uint64(d.UserID),
		PaymentMethod:  sql.NullInt64{Int64: int64(d.PaymentMethod), Valid: true},
		OriginalAmount: d.Amounts.OriginAmount,
		DiscountAmount: d.Amounts.DiscountAmount,
		PayableAmount:  d.Amounts.FinalAmount,
		OrderStatus:    int64(order.OrderStatus_ORDER_STATUS_CREATED),
		PaymentStatus:  int64(order.PaymentStatus_PAYMENT_STATUS_NOT_PAID),
		ExpireTime:     time.Now().Add(biz.OrderExpireTime).Unix(),
	}
}
func (l *CreateOrderLogic) generateOrderID() string {
	var uid uuid.UUID
	uid, err := uuid.NewV7()
	if err != nil {
		l.Logger.Infow("uuid generate failed", logx.Field("err", err))
		uid = uuid.New()
	}
	OrderID := uid.String()
	return OrderID
}
func (l *CreateOrderLogic) logContext(dto *orderCreateDTO) []logx.LogField {
	return []logx.LogField{
		logx.Field("user_id", dto.UserID),
		logx.Field("pre_order_id", dto.PreOrderID),
		logx.Field("order_id", dto.OrderID),
	}
}
