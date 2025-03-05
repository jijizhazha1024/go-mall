package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	var queryErr error
	paymentModel := l.svcCtx.PaymentModel

	// 查询商品列表
	offset := (in.Pagination.Page - 1) * in.Pagination.PageSize
	payments, queryErr := paymentModel.FindPage(l.ctx, in.UserId, int(offset), int(in.Pagination.PageSize))
	// 统一错误处理
	if queryErr != nil {
		if errors.Is(queryErr, sqlx.ErrNotFound) {
			return &payment.PaymentListResp{}, nil
		}
		l.Logger.Errorw("query payments failed",
			logx.Field("err", queryErr),
			logx.Field("page", in.Pagination.Page),
			logx.Field("pageSize", in.Pagination.PageSize))
		return nil, queryErr
	}
	items := make([]*payment.PaymentItem, len(payments))
	for i, p := range payments {
		items[i] = &payment.PaymentItem{
			PaymentId:      p.PaymentId,
			PreOrderId:     p.PreOrderId,
			OrderId:        p.OrderId.String,
			OriginalAmount: p.OriginalAmount,
			PaidAmount:     p.PaidAmount.Int64,
			TransactionId:  p.TransactionId.String,
			PayUrl:         p.PayUrl,
			ExpireTime:     p.ExpireTime,
			Status:         payment.PaymentStatus(p.Status),
			CreatedAt:      p.CreatedAt.Unix(),
			UpdatedAt:      p.UpdatedAt.Unix(),
			PaidAt:         p.PaidAt.Int64,
		}
	}

	return &payment.PaymentListResp{
		Payments:   items,
	}, nil
}
