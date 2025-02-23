package logic

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	inventory2 "jijizhazha1024/go-mall/dal/model/inventory/inventory"
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

	for _, item := range in.Items {

		if item.Quantity <= 0 {
			l.Logger.Infow("quantity must be greater than 0", logx.Field("quantity", item.Quantity), logx.Field("product_id", item.ProductId))
			return nil, biz.InvalidInventoryErr
		}
		tostr := fmt.Sprintf("%d", item.Quantity)

		err := l.svcCtx.Rdb.Set(fmt.Sprintf("inventory:product:%d", item.ProductId), tostr)
		if err != nil {
			l.Logger.Errorw("update inventory failed", logx.Field("product_id", item.ProductId), logx.Field("err", err))
			return nil, err
		}
		//执行sql
		if err := l.svcCtx.InventoryModel.UpdateOrCreate(l.ctx, inventory2.Inventory{
			ProductId: int64(item.ProductId),
			Total:     int64(item.Quantity),
		}); err != nil {
			l.Logger.Errorw("update inventory error", logx.Field("error", err.Error()), logx.Field("product_id", item.ProductId))
			return nil, err
		}
	}
	return &inventory.InventoryResp{}, nil

}
