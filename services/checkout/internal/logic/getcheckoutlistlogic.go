package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCheckoutListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCheckoutListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCheckoutListLogic {
	return &GetCheckoutListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCheckoutListLogic) GetCheckoutList(in *checkout.CheckoutListReq) (*checkout.CheckoutListResp, error) {
	// todo: add your logic here and delete this line

	return &checkout.CheckoutListResp{}, nil
}
