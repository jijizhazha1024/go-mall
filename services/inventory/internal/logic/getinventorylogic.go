package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"

	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryLogic {
	return &GetInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetInventory 查询库存
func (l *GetInventoryLogic) GetInventory(in *inventory.GetInventoryReq) (*inventory.GetInventoryResp, error) {
	inventoryResp, err := l.svcCtx.InventoryModel.FindOne(l.ctx, int64(in.ProductId))
	res := new(inventory.GetInventoryResp)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			l.Logger.Infow("product not in inventory", logx.Field("product_id", in.ProductId))
			res.StatusCode = code.ProductNotFoundInventory
			res.StatusMsg = code.ProductNotFoundInventoryMsg
			return res, nil
		}
		l.Logger.Errorw("product inventory get failed", logx.Field("product_id", in.ProductId))
		return nil, err
	}
	inventoryCnt := inventoryResp.Total - inventoryResp.Sold
	res.SoldCount = inventoryResp.Total
	// 存在库存
	if inventoryCnt > 0 {
		res.Inventory = inventoryCnt
		res.SoldCount = inventoryResp.Sold
	}
	return res, nil
}
