package logic

import (
	"context"
	"jijizhazha1024/go-mall/common/consts/code"
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

// GetCheckoutList 获取结算列表
func (l *GetCheckoutListLogic) GetCheckoutList(in *checkout.CheckoutListReq) (*checkout.CheckoutListResp, error) {
	// 1. 参数校验
	if in.Page < 1 {
		in.Page = 1
	}
	if in.PageSize <= 0 || in.PageSize > 50 {
		in.PageSize = 10
	}

	// 2. 查询总记录数
	total, err := l.svcCtx.CheckoutModel.CountByUserId(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Errorw("查询结算订单总数失败",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId))
		return &checkout.CheckoutListResp{
			StatusCode: code.QueryOrderTotalFailed,
			StatusMsg:  code.QueryOrderTotalFailedMsg,
		}, nil
	}

	// 3. 如果没有记录，直接返回空列表
	if total == 0 {
		return &checkout.CheckoutListResp{
			Total: 0,
			Data:  []*checkout.CheckoutOrder{},
		}, nil
	}

	// 4. 查询分页订单数据
	checkouts, err := l.svcCtx.CheckoutModel.FindByUserId(l.ctx, in.UserId, in.Page, in.PageSize)
	if err != nil {
		l.Logger.Errorw("查询结算列表失败",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId),
			logx.Field("page", in.Page),
			logx.Field("page_size", in.GetPageSize()))
		return &checkout.CheckoutListResp{
			StatusCode: code.QueryOrderListFailed,
			StatusMsg:  code.QueryOrderListFailedMsg,
		}, nil
	}
	// 5. 查询所有订单对应的商品详情
	preOrderIds := make([]string, 0, len(checkouts))
	for _, c := range checkouts {
		preOrderIds = append(preOrderIds, c.PreOrderId)
	}

	itemsMap, err := l.svcCtx.CheckoutItemsModel.FindItemsByPreOrderIds(l.ctx, preOrderIds)
	if err != nil {
		l.Logger.Errorw("查询订单商品详情失败",
			logx.Field("err", err),
			logx.Field("pre_order_ids", preOrderIds))
		return &checkout.CheckoutListResp{
			StatusCode: code.QueryOrderProductFailed,
			StatusMsg:  code.QueryOrderProductFailedMsg,
		}, nil
	}

	// 6. 组装数据
	var checkoutOrders []*checkout.CheckoutOrder
	for _, c := range checkouts {
		order := &checkout.CheckoutOrder{
			PreOrderId:     c.PreOrderId,
			UserId:         int64(c.UserId),
			Status:         checkout.CheckoutStatus(c.Status),
			ExpireTime:     c.ExpireTime,
			CreatedAt:      c.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      c.UpdatedAt.Format("2006-01-02 15:04:05"),
			Items:          itemsMap[c.PreOrderId],
			OriginalAmount: c.OriginalAmount,
			FinalAmount:    c.FinalAmount,
		}
		checkoutOrders = append(checkoutOrders, order)
	}

	// 7. 返回数据
	return &checkout.CheckoutListResp{
		Total: total,
		Data:  checkoutOrders,
	}, nil
}
