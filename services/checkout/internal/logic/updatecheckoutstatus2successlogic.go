package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCheckoutStatus2SuccessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCheckoutStatus2SuccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCheckoutStatus2SuccessLogic {
	return &UpdateCheckoutStatus2SuccessLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCheckoutStatus2SuccessLogic) UpdateCheckoutStatus2Success(in *checkout.UpdateCheckoutStatusReq) (*checkout.UpdateCheckoutStatusResp, error) {
	// todo: add your logic here and delete this line

	return &checkout.UpdateCheckoutStatusResp{}, nil
}
