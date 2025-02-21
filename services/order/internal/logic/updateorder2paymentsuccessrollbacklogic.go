package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrder2PaymentSuccessRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrder2PaymentSuccessRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrder2PaymentSuccessRollbackLogic {
	return &UpdateOrder2PaymentSuccessRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrder2PaymentSuccessRollback 支付失败的补充操作
func (l *UpdateOrder2PaymentSuccessRollbackLogic) UpdateOrder2PaymentSuccessRollback(in *order.UpdateOrder2PaymentSuccessRequest) (*order.EmptyRes, error) {
	// todo: add your logic here and delete this line

	return &order.EmptyRes{}, nil
}
