package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/common/consts/code"
	order2 "jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/coupons/coupons"
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
func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderRequest) (*order.OrderResponse, error) {
	// --------------- 参数校验 ---------------
	if err := l.validateRequest(in); err != nil {
		return nil, err
	}

	// --------------- 绑定数据 ---------------
	dto, err := l.collectOrderData(in)
	if err != nil {
		l.Logger.Errorw("collect order data failed", logx.Field("user_id", in.UserId),
			logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err))
		return nil, err
	}

	var uid uuid.UUID
	uid, err = uuid.NewV7()
	if err != nil {
		uid = uuid.New()
	}
	OrderID := uid.String()

	res := &order.OrderResponse{}
	expireTime := time.Minute * 30
	orderValue := &order2.Orders{
		OrderId:        OrderID,
		PreOrderId:     in.PreOrderId,
		UserId:         uint64(in.UserId),
		PaymentMethod:  sql.NullInt64{Int64: int64(in.PaymentMethod), Valid: true},
		OriginalAmount: dto.Amounts.OriginAmount,
		DiscountAmount: dto.Amounts.DiscountAmount,
		PayableAmount:  dto.Amounts.FinalAmount,
		OrderStatus:    int64(order.OrderStatus_ORDER_STATUS_CREATED),
		PaymentStatus:  int64(order.PaymentStatus_PAYMENT_STATUS_NOT_PAID),
		ExpireTime:     time.Now().Add(expireTime).Unix(),
	}

	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {

		// --------------- check ---------------
		// 订单是否存在，通过 pre_order_id 查询
		orderRes, err := l.svcCtx.OrderModel.WithSession(session).CheckOrderExistByPreOrderId(ctx, in.PreOrderId, int32(in.UserId))
		if err != nil {
			l.Logger.Errorw("check order exist failed", logx.Field("user_id", in.UserId),
				logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err))
			return err
		}
		if orderRes {
			l.Logger.Infow("transaction aborted", logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
			res.StatusCode = code.OrderExist
			res.StatusMsg = code.OrderExistMsg
			return nil
		}

		// --------------- insert ---------------
		// 插入订单
		_, err = l.svcCtx.OrderModel.WithSession(session).Insert(l.ctx, orderValue)
		if err != nil {
			l.Logger.Errorw("insert order failed", logx.Field("user_id", in.UserId),
				logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err))
			return err
		}
		// 插入订单项关联商品
		for _, item := range dto.Items {
			orderItem := &order2.OrderItems{
				OrderId:   OrderID,
				ProductId: uint64(item.ProductId),
				Quantity:  uint64(item.Quantity),
			}
			if _, err = l.svcCtx.OrderItemModel.WithSession(session).Insert(l.ctx, orderItem); err != nil {
				l.Logger.Errorw("insert order item failed", logx.Field("user_id", in.UserId),
					logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err))
				return err
			}
		}
		// 插入订单地址
		dto.Address.OrderId = OrderID
		if _, err := l.svcCtx.OrderAddress.Insert(l.ctx, dto.Address); err != nil {
			l.Logger.Errorw("insert order address failed", logx.Field("user_id", in.UserId),
				logx.Field("pre_order_id", in.PreOrderId),
			)
		}
		return nil
	}); err != nil {
		l.Logger.Errorw("insert order failed", logx.Field("user_id", in.UserId),
			logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err))
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	if res.StatusCode != code.Success {
		l.Logger.Infow("transaction aborted", logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
		return nil, status.Error(codes.Aborted, res.StatusMsg)
	}
	return &order.OrderResponse{}, nil
}
func (l *CreateOrderLogic) validateRequest(in *order.CreateOrderRequest) error {
	if in.PreOrderId == "" || in.UserId == 0 || in.AddressId == 0 || in.CouponId == "" || in.PaymentMethod == 0 {
		return status.Error(codes.InvalidArgument, "参数不合法")
	}
	return nil
}
func (l *CreateOrderLogic) collectOrderData(in *order.CreateOrderRequest) (*orderCreateDTO, error) {
	g, ctx := errgroup.WithContext(l.ctx)
	var dto = &orderCreateDTO{
		PreOrderID:    in.PreOrderId,
		UserID:        int64(in.UserId),
		PaymentMethod: in.PaymentMethod,
	}
	g.Go(func() error {
		// 获取订单详情
		checkoutDetail, err := l.svcCtx.CheckoutRpc.GetCheckoutDetail(ctx, &checkout.CheckoutDetailReq{
			PreOrderId: in.PreOrderId,
		})
		if err != nil {
			logx.Errorw("call rpc GetCheckoutDetail failed", logx.Field("user_id", in.UserId),
				logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err))
			return err
		}
		if checkoutDetail.StatusCode != code.Success {
			logx.Errorw("call rpc GetCheckoutDetail failed", logx.Field("user_id", in.UserId),
				logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err),
			)
			return status.Error(codes.Aborted, checkoutDetail.StatusMsg)

		}
		// 计算优惠价格
		couponResp, err := l.svcCtx.CouponRpc.CalculateCoupon(ctx, &coupons.CalculateCouponReq{
			UserId:   int32(in.UserId),
			CouponId: in.CouponId,
			Items:    convertToCouponItems(checkoutDetail.Data.Items),
		})
		if err != nil {
			logx.Errorw("call rpc CalculateCoupon failed", logx.Field("user_id", in.UserId),
				logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err),
			)
			return err
		}
		if couponResp.StatusCode != code.Success {
			logx.Errorw("call rpc CalculateCoupon failed", logx.Field("user_id", in.UserId),
				logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err),
			)
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
			l.Logger.Errorw("call rpc GetAddress failed", logx.Field("user_id", in.UserId),
				logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err))
			return err
		}
		if addressResp.StatusCode != code.Success {
			l.Logger.Infow("transaction aborted", logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
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
func convertToCouponItems(items []*checkout.CheckoutItem) []*coupons.Items {
	couponItems := make([]*coupons.Items, 0, len(items))
	for _, item := range items {
		couponItems = append(couponItems, &coupons.Items{
			ProductId: int32(item.ProductId),
			Quantity:  item.Quantity,
		})
	}
	return couponItems
}
