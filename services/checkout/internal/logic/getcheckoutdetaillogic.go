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
	checkoutRecord, err := l.svcCtx.CheckoutModel.FindOneByUserIdAndPreOrderId(l.ctx, in.UserId, in.PreOrderId)
	if err != nil {
		if err != nil {
			return nil, err
		}
		l.Logger.Errorw("查询结算记录失败", logx.Field("err", err), logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
		return &checkout.CheckoutDetailResp{
			StatusCode: 500,
			StatusMsg:  "查询结算记录失败",
		}, err
	}

	checkoutItems, err := l.svcCtx.CheckoutItemsModel.FindOne(l.ctx, in.PreOrderId)
	if err != nil {
		if err != nil {
			return nil, err
		}
		l.Logger.Errorw("查询结算记录失败", logx.Field("err", err), logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
		return &checkout.CheckoutDetailResp{
			StatusCode: 500,
			StatusMsg:  "查询结算记录失败",
		}, err
	}

	orderData := &checkout.CheckoutOrder{
		PreOrderId: checkoutRecord.PreOrderId,
		UserId:     int64(checkoutRecord.UserId),
		Status:     checkout.CheckoutStatus(checkoutRecord.Status),
		ExpireTime: checkoutRecord.ExpireTime,
		CreatedAt:  checkoutRecord.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  checkoutRecord.UpdatedAt.String(),
	}

	var items []*checkout.CheckoutItem
	if checkoutItems != nil {
		checkoutItem := &checkout.CheckoutItem{
			ProductId: int32(checkoutItems.ProductId),
			Quantity:  int32(checkoutItems.Quantity),
			Price:     checkoutItems.Price,
		}
		items = append(items, checkoutItem)
	}

	orderData.Items = items

	resp := &checkout.CheckoutDetailResp{
		StatusCode: 200,
		StatusMsg:  "成功获取结算详情",
		Data:       orderData,
	}

	return resp, nil
}
