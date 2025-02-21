package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrder2PaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrder2PaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrder2PaymentLogic {
	return &UpdateOrder2PaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrder2Payment 更新订单（支付服务回调使用）
func (l *UpdateOrder2PaymentLogic) UpdateOrder2Payment(in *order.UpdateOrder2PaymentRequest) (*order.EmptyRes, error) {
	// todo: add your logic here and delete this line

	return &order.EmptyRes{}, nil
}
