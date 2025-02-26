package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/errgroup"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrder 获取订单详情
func (l *GetOrderLogic) GetOrder(in *order.GetOrderRequest) (*order.OrderDetailResponse, error) {
	res := &order.OrderDetailResponse{}
	if in.OrderId == "" || in.UserId == 0 {
		res.StatusCode = 400
		return res, nil
	}
	group, ctx := errgroup.WithContext(l.ctx)
	// mapreduce
	// 获取订单
	group.Go(func() error {
		one, err := l.svcCtx.OrderModel.GetOrderByOrderIDAndUserIDWithLock(ctx, in.OrderId, int32(in.UserId))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				res.StatusCode = code.OrderNotExist
				res.StatusMsg = code.OrderNotExistMsg
				return err
			}
			logx.Errorw("query order failed", logx.Field("err", err))
			return err
		}
		res.Order = convertToOrderResp(one)
		return nil
	})
	// 获取订单关联商品
	group.Go(func() error {
		items, err := l.svcCtx.OrderItemModel.QueryOrderItemsByOrderID(ctx, in.OrderId)
		if err != nil {
			logx.Errorw("query order items failed", logx.Field("err", err))
			return err
		}
		if len(items) == 0 {
			res.StatusCode = code.UserOrderItemNotExist
			res.StatusMsg = code.UserOrderItemNotExistMsg
			return err
		}
		res.Items = convertToOrderItemResp(items)
		return nil
	})
	// 获取订单关联地址
	group.Go(func() error {
		address, err := l.svcCtx.OrderAddress.GetOrderAddressByOrderID(ctx, in.OrderId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				res.StatusCode = code.UserOrderAddressNotExist
				res.StatusMsg = code.UserOrderAddressNotExistMsg
				return err
			}
			logx.Errorw("query order address failed", logx.Field("err", err))
			return err
		}
		res.Address = convertToOrderAddressResp(address)
		return nil
	})
	if err := group.Wait(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			res.StatusCode = code.OrderNotExist
			res.StatusMsg = code.OrderNotExistMsg
			return res, nil
		}
		return nil, err
	}
	return res, nil
}
