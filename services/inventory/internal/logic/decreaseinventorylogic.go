package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"

	"github.com/zeromicro/go-zero/core/logx"
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

	if in.Quantity <= 0 {
		l.Logger.Infow("quantity must be greater than 0", logx.Field("quantity", in.Quantity), logx.Field("product_id", in.ProductId))
		return nil, biz.InvalidInventoryErr
	}
	// 事务
	cnt, err := l.svcCtx.InventoryModel.DecreaseInventoryAtom(l.ctx, in.ProductId, in.Quantity)

	switch {
	case errors.Is(err, sqlx.ErrNotFound):
		l.Logger.Infow("product not in inventory", logx.Field("product_id", in.ProductId))
		res.StatusCode = code.ProductNotFoundInventory
		res.StatusMsg = code.ProductNotFoundInventoryMsg
		return res, nil

	case errors.Is(err, biz.InventoryNotEnoughErr):
		l.Logger.Infow("product inventory not enough", logx.Field("product_id", in.ProductId))
		res.StatusCode = code.InventoryNotEnough
		res.StatusMsg = code.InventoryNotEnoughMsg
		return res, nil
	case errors.Is(err, biz.InventoryDecreaseFailedErr):
		l.Logger.Errorw("product inventory decrease failed", logx.Field("product_id", in.ProductId))
		return nil, err
	}
	res.Inventory = cnt
	return res, nil
}
