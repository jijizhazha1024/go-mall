package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"
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
	err := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		checkoutRecord, err := l.svcCtx.CheckoutModel.FindOneByUserIdAndPreOrderIdWithSession(l.ctx, session, in.UserId, in.PreOrderId)
		if err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return status.Error(codes.Aborted, code.OutOfRecordMsg)
			} else {
				return err
			}
		}
		switch checkout.CheckoutStatus(checkoutRecord.Status) {
		case checkout.CheckoutStatus_CONFIRMED:
			l.Logger.Infof("订单 %s 已经是已确认状态", in.PreOrderId)
			return nil

		case checkout.CheckoutStatus_CANCELLED, checkout.CheckoutStatus_EXPIRED:
			// 订单已经过期进行回滚
			return status.Error(codes.Aborted, code.OrderStatusInvalidMsg)
		}

		err = l.svcCtx.CheckoutModel.UpdateStatusWithSession(l.ctx, session, int64(checkout.CheckoutStatus_CONFIRMED), in.UserId, in.PreOrderId)
		if err != nil {
			return err
		}

		l.Logger.Infof("成功更新订单 %s 的结算状态为已确认", in.PreOrderId)
		return nil
	})

	if err != nil {
		l.Logger.Errorw("事务处理失败",
			logx.Field("err", err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &checkout.EmptyResp{}, nil
}
