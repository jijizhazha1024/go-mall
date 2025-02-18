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

type ReturnInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReturnInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReturnInventoryLogic {
	return &ReturnInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReturnInventory 归还库存
func (l *ReturnInventoryLogic) ReturnInventory(in *inventory.InventoryReq) (*inventory.InventoryResp, error) {
	var res = new(inventory.InventoryResp)

	if in.Quantity <= 0 {
		l.Logger.Infow("quantity must be greater than 0", logx.Field("quantity", in.Quantity), logx.Field("product_id", in.ProductId))
		return nil, biz.InvalidInventoryErr
	}
	cnt, err := l.svcCtx.InventoryModel.ReturnInventory(l.ctx, in.ProductId, in.Quantity)
	switch {
	case errors.Is(err, sqlx.ErrNotFound):
		l.Logger.Infow("product not in inventory", logx.Field("product_id", in.ProductId))
		res.StatusCode = code.ProductNotFoundInventory
		res.StatusMsg = code.ProductNotFoundInventoryMsg
		return res, nil

	case errors.Is(err, biz.InventoryDecreaseFailedErr):
		l.Logger.Errorw("product inventory decrease failed", logx.Field("product_id", in.ProductId))
		return nil, err

	}
	res.Inventory = cnt
	return res, nil
}
