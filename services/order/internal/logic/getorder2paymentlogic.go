package logic

import (
	"context"
	"database/sql"
	"errors"
	"jijizhazha1024/go-mall/common/consts/code"

	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrder2PaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrder2PaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrder2PaymentLogic {
	return &GetOrder2PaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrder2PaymentLogic) GetOrder2Payment(in *order.GetOrderRequest) (*order.OrderDetail2PaymentResponse, error) {
	res := &order.OrderDetail2PaymentResponse{}
	if in.OrderId == "" || in.UserId == 0 {
		res.StatusCode = code.OrderParameterInvalid
		res.StatusMsg = code.OrderParameterInvalidMsg
		return res, nil
	}
	one, err := l.svcCtx.OrderModel.GetOrderByOrderIDAndUserIDWithLock(l.ctx, in.OrderId, int32(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res.StatusCode = code.OrderNotExist
			res.StatusMsg = code.OrderNotExistMsg
			return res, nil
		}
		logx.Errorw("query order failed", logx.Field("err", err))
		return nil, err
	}
	res.Order = convertToOrderResp(one)
	return res, nil
}
