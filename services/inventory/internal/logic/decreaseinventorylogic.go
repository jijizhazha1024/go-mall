package logic

import (
	"context"
	"errors"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DecreaseInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecreaseInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecreaseInventoryLogic {
	return &DecreaseInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DecreaseInventory 扣减库存
func (l *DecreaseInventoryLogic) DecreaseInventory(in *inventory.InventoryReq) (*inventory.InventoryResp, error) {

	var res = new(inventory.InventoryResp)

	//将id和数量分别存入数组
	productId := make([]int32, 0)
	quantity := make([]int32, 0)
	for _, item := range in.Items {
		productId = append(productId, item.ProductId)
		quantity = append(quantity, item.Quantity)
	}

	// 事务
	err := l.svcCtx.InventoryModel.BatchDecreaseInventoryAtom(l.ctx, productId, quantity)

	switch {
	case errors.Is(err, sqlx.ErrNotFound):
		l.Logger.Infow("product not in inventory", logx.Field("product_id", productId))
		res.StatusCode = code.ProductNotFoundInventory
		res.StatusMsg = code.ProductNotFoundInventoryMsg
		return res, nil

	case errors.Is(err, biz.InventoryNotEnoughErr):
		l.Logger.Infow("product inventory not enough", logx.Field("product_id", productId))
		res.StatusCode = code.InventoryNotEnough
		res.StatusMsg = code.InventoryNotEnoughMsg
		return res, nil

	case errors.Is(err, biz.InventoryDecreaseFailedErr):
		l.Logger.Errorw("product inventory decrease failed", logx.Field("product_id", productId))
		return nil, err
	}
	if err != nil {
		l.Logger.Errorw("product inventory decrease failed", logx.Field("product_id", productId))
		return nil, err
	}

	return res, nil
}
