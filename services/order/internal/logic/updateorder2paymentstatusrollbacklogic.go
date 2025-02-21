package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrder2PaymentStatusRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrder2PaymentStatusRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrder2PaymentStatusRollbackLogic {
	return &UpdateOrder2PaymentStatusRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrder2PaymentStatusRollback 补偿操作 更新订单（支付服务回调使用） 创建状态
func (l *UpdateOrder2PaymentStatusRollbackLogic) UpdateOrder2PaymentStatusRollback(in *order.UpdateOrder2PaymentRequest) (*order.EmptyRes, error) {
	// todo: add your logic here and delete this line

	return &order.EmptyRes{}, nil
}
