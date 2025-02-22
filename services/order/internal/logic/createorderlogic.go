package logic

import (
	"context"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/common/consts/code"
	order2 "jijizhazha1024/go-mall/dal/model/order"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/coupons/coupons"

	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

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
	// --------------- 获取结算信息 ---------------
	// 此时结算状态为已确认。
	// 获取结算信息
	checkoutDetail, err := l.svcCtx.CheckoutRpc.GetCheckoutDetail(l.ctx, &checkout.CheckoutDetailReq{
		PreOrderId: in.PreOrderId,
	})
	if err != nil {
		l.Logger.Errorw("call GetCheckoutDetail rpc failed", logx.Field("user_id", in.UserId),
			logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err))
		// 进行回滚
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	if checkoutDetail.StatusCode != code.Success {
		l.Logger.Infow("transaction aborted", logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
		// 进行回滚
		return nil, status.Error(codes.Aborted, checkoutDetail.StatusMsg)
	}
	// 校验结算状态，确认状态，在创建订单时先会调用结算服务，说以这里只需要校验确认状态，无需确保结算超时
	if checkoutDetail.Data.Status != checkout.CheckoutStatus_CONFIRMED {
		l.Logger.Infow("transaction aborted", logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
		return nil, status.Error(codes.Aborted, "结算状态异常")
	}

	// --------------- 获取优惠券 ---------------
	items := make([]*coupons.Items, 0, len(checkoutDetail.Data.Items))
	for _, item := range checkoutDetail.Data.Items {
		items = append(items, &coupons.Items{
			ProductId: int32(item.ProductId),
			Quantity:  item.Quantity,
		})
	}
	// 此时优惠券已经是存在的，且以处于预锁定状态
	couponResp, err := l.svcCtx.CouponRpc.CalculateCoupon(l.ctx, &coupons.CalculateCouponReq{
		UserId:   int32(in.UserId),
		CouponId: in.CouponId,
		Items:    items,
	})
	if err != nil {
		l.Logger.Errorw("call rpc CalculateCoupon failed", logx.Field("user_id", in.UserId),
			logx.Field("pre_order_id", in.PreOrderId), logx.Field("err", err))
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	if couponResp.StatusCode != code.Success {
		l.Logger.Infow("transaction aborted", logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
		return nil, status.Error(codes.Aborted, couponResp.StatusMsg)
	}

	res := &order.OrderResponse{}
	//orderInserValue := &order2.Orders{
	//	PreOrderId:     in.PreOrderId,
	//	UserId:         uint64(in.UserId),
	//	PaymentMethod:  sql.NullInt64{Int64: int64(in.PaymentMethod), Valid: true},
	//	//OriginalAmount: couponResp.OriginAmount,
	//}
	if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
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
		l.svcCtx.OrderModel.WithSession(session).Insert(l.ctx, &order2.Orders{})

		return nil
	}); err != nil {

	}
	// 由消息队列去创建商品信息快照

	return &order.OrderResponse{}, nil
}
