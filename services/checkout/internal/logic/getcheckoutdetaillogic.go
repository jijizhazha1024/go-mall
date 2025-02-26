package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCheckoutDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCheckoutDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCheckoutDetailLogic {
	return &GetCheckoutDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCheckoutDetail 获取结算详情
func (l *GetCheckoutDetailLogic) GetCheckoutDetail(in *checkout.CheckoutDetailReq) (*checkout.CheckoutDetailResp, error) {
	// todo: add your logic here and delete this line

	return &checkout.CheckoutDetailResp{}, nil
}
