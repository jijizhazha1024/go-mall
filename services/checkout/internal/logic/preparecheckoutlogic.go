package logic

import (
	"context"
	"crypto/md5"
	"fmt"

	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PrepareCheckoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPrepareCheckoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareCheckoutLogic {
	return &PrepareCheckoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 生成用于幂等性的唯一哈希
func generateHash(productId int32, quantity int32) string {
	hashString := fmt.Sprintf("%d:%d", productId, quantity)
	hashBytes := md5.Sum([]byte(hashString))
	return fmt.Sprintf("%x", hashBytes)
}

// PrepareCheckout 预结算)生成预订单）
func (l *PrepareCheckoutLogic) PrepareCheckout(in *checkout.CheckoutReq) (*checkout.CheckoutResp, error) {
	// todo: add your logic here and delete this line

	return &checkout.CheckoutResp{}, nil
}
