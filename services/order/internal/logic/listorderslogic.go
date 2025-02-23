package logic

import (
	"context"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/order/internal/svc"
	"jijizhazha1024/go-mall/services/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListOrders 分页查询订单列表
func (l *ListOrdersLogic) ListOrders(in *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {

	res := &order.ListOrdersResponse{}
	// --------------- check ---------------
	if in.UserId == 0 {
		res.StatusCode = code.UserNotFound
		res.StatusMsg = code.UserNotFoundMsg
		return res, nil
	}
	if in.Pagination.PageSize <= 0 || in.Pagination.PageSize > biz.MaxPageSize {
		in.Pagination.PageSize = biz.MaxPageSize
	}
	if in.Pagination.Page <= 0 {
		in.Pagination.Page = 1
	}
	orderList, err := l.svcCtx.OrderModel.GetOrdersByUserID(l.ctx, int32(in.UserId), in.Pagination.Page, in.Pagination.PageSize)
	if err != nil {
		l.Logger.Errorw("call svcCtx.OrderModel.GetOrdersByUserID failed", logx.Field("err", err))
		res.StatusCode = code.ServerError
		res.StatusMsg = code.ServerErrorMsg
		return res, nil
	}
	res.Orders = make([]*order.Order, len(orderList))
	for i, o := range orderList {
		res.Orders[i] = convertToOrderResp(o)
	}
	return res, nil
}
