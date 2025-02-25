package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReleaseCheckoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReleaseCheckoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseCheckoutLogic {
	return &ReleaseCheckoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateCheckoutStatus2Success 当订单超时，支付超时，支付退款
func (l *ReleaseCheckoutLogic) ReleaseCheckout(in *checkout.ReleaseReq) (*checkout.EmptyResp, error) {
	// todo: add your logic here and delete this line

	return &checkout.EmptyResp{}, nil
}
