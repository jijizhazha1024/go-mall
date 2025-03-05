package logic

import (
	"context"
	xerrors "github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/checkout/checkout"

	"jijizhazha1024/go-mall/apis/checkout/internal/svc"
	"jijizhazha1024/go-mall/apis/checkout/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckoutLogic {
	return &CheckoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckoutLogic) Checkout(req *types.CheckoutReq) (resp *types.CheckoutResp, err error) {
	userID, ok := l.ctx.Value(biz.UserIDKey).(uint32)
	if !ok {
		return nil, xerrors.New(code.AuthBlank, code.AuthBlankMsg)
	}
	res, err := l.svcCtx.CheckoutRpc.PrepareCheckout(l.ctx, &checkout.CheckoutReq{
		UserId:     userID,
		CouponId:   req.CouponID,
		OrderItems: convertCheckoutItem2Req(req.OrderItems),
	})
	if err != nil {
		l.Logger.Errorw("call rpc GetOrder failed", logx.Field("err", err))
		return nil, xerrors.New(code.ServerError, code.ServerErrorMsg)
	}
	if res.StatusCode != code.Success {
		return nil, xerrors.New(int(res.StatusCode), res.StatusMsg)
	}
	resp = &types.CheckoutResp{
		ExpireTime: res.ExpireTime,
		PayMethod:  res.PayMethod,
		PreOrderID: res.PreOrderId,
	}
	return
}

func convertCheckoutItem2Req(items []types.CheckoutItemReq) []*checkout.CheckoutReq_OrderItem {
	orderItems := make([]*checkout.CheckoutReq_OrderItem, len(items))
	for i, item := range items {
		orderItems[i] = &checkout.CheckoutReq_OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	return orderItems
}
