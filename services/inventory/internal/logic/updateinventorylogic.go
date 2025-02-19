package logic

import (
	"context"
	"jijizhazha1024/go-mall/common/consts/biz"
	inventory2 "jijizhazha1024/go-mall/dal/model/inventory"
	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInventoryLogic {
	return &UpdateInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateInventory 更新库存，进行修改库存数量
func (l *UpdateInventoryLogic) UpdateInventory(in *inventory.InventoryReq) (*inventory.InventoryResp, error) {

	if in.Quantity <= 0 {
		l.Logger.Infow("quantity must be greater than 0", logx.Field("quantity", in.Quantity), logx.Field("product_id", in.ProductId))
		return nil, biz.InvalidInventoryErr
	}
	if err := l.svcCtx.InventoryModel.UpdateOrCreate(l.ctx, inventory2.Inventory{
		ProductId: int64(in.ProductId),
		Total:     int64(in.Quantity),
	}); err != nil {
		l.Logger.Errorw("update inventory error", logx.Field("error", err.Error()), logx.Field("product_id", in.ProductId))
		return nil, err
	}
	return &inventory.InventoryResp{}, nil
}
