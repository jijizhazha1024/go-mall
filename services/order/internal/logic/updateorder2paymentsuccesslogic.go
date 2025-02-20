package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrder2PaymentSuccessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrder2PaymentSuccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrder2PaymentSuccessLogic {
	return &UpdateOrder2PaymentSuccessLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 支付成功时（进行修改订单状态）
func (l *UpdateOrder2PaymentSuccessLogic) UpdateOrder2PaymentSuccess(in *order.UpdateOrder2PaymentSuccessRequest) (*order.EmptyRes, error) {
	// todo: add your logic here and delete this line

	return &order.EmptyRes{}, nil
}
