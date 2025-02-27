package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"log"

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
		if errors.Is(err, sqlx.ErrNotFound) {
			log.Println("No items found for the given userId and preOrderId.")
		} else {
			l.Logger.Errorw("查询结算记录失败", logx.Field("err", err), logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
			return nil, err
		}
	}

	checkoutItems, err := l.svcCtx.CheckoutItemsModel.FindItemsByUserAndPreOrder(l.ctx, in.UserId, in.PreOrderId)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			log.Println("No items found for the given userId and preOrderId.")
		} else {
			l.Logger.Errorw("查询结算记录失败", logx.Field("err", err), logx.Field("user_id", in.UserId), logx.Field("pre_order_id", in.PreOrderId))
			return nil, err
		}
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
	for _, item := range checkoutItems {
		checkoutItem := &checkout.CheckoutItem{
			ProductId: int32(item.ProductId),
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		}
		items = append(items, checkoutItem)
	}

	orderData.Items = items

	resp := &checkout.CheckoutDetailResp{
		Data: orderData,
	}

	return resp, nil
}
