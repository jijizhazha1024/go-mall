package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStatus2OrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStatus2OrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStatus2OrderLogic {
	return &UpdateStatus2OrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateStatus2Order 由订单服务调用，更新结算状态为已确认
func (l *UpdateStatus2OrderLogic) UpdateStatus2Order(in *checkout.UpdateStatusReq) (*checkout.EmptyResp, error) {
	// todo: add your logic here and delete this line

	return &checkout.EmptyResp{}, nil
}
