package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrder2PaymentStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrder2PaymentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrder2PaymentStatusLogic {
	return &UpdateOrder2PaymentStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrder2Payment 更新订单（支付服务回调使用） 更新为支付中
func (l *UpdateOrder2PaymentStatusLogic) UpdateOrder2PaymentStatus(in *order.UpdateOrder2PaymentRequest) (*order.EmptyRes, error) {
	// todo: add your logic here and delete this line

	return &order.EmptyRes{}, nil
}
