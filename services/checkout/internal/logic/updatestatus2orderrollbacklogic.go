package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStatus2OrderRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStatus2OrderRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStatus2OrderRollbackLogic {
	return &UpdateStatus2OrderRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateStatus2OrderRollback 补偿操作
func (l *UpdateStatus2OrderRollbackLogic) UpdateStatus2OrderRollback(in *checkout.UpdateStatusReq) (*checkout.EmptyResp, error) {
	// todo: add your logic here and delete this line

	return &checkout.EmptyResp{}, nil
}
