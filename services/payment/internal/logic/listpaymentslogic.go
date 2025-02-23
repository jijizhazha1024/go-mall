package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/payment/internal/svc"
	"jijizhazha1024/go-mall/services/payment/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPaymentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPaymentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPaymentsLogic {
	return &ListPaymentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPaymentsLogic) ListPayments(in *payment.PaymentListReq) (*payment.PaymentListResp, error) {
	// todo: add your logic here and delete this line

	return &payment.PaymentListResp{}, nil
}
