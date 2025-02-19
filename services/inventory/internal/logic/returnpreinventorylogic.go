package logic

import (
	"context"
	"errors"
	"fmt"

	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ReturnPreInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReturnPreInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReturnPreInventoryLogic {
	return &ReturnPreInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReturnPreInventory 退还预扣减的库存（）
func (l *ReturnPreInventoryLogic) ReturnPreInventory(in *inventory.InventoryReq) (*inventory.InventoryResp, error) {

	var res = new(inventory.InventoryResp)

	if in.Quantity <= 0 {
		l.Logger.Infow("quantity must be greater than 0", logx.Field("quantity", in.Quantity), logx.Field("product_id", in.ProductId))
		return nil, biz.InvalidInventoryErr
	}
	// 获取缓存内的总量
	total, err := l.svcCtx.Rdb.Hget(fmt.Sprintf("inventory:%d", in.ProductId), "total")

	if err != nil {
		l.Logger.Errorw("get inventory total failed", logx.Field("product_id", in.ProductId), logx.Field("err", err))
		return nil, err
	}

	if total == "" {
		l.Logger.Infow("product not in inventory", logx.Field("product_id", in.ProductId))
		res.StatusCode = code.ProductNotFoundInventory
		res.StatusMsg = code.ProductNotFoundInventoryMsg
		return res, nil
	}

	// 增库存
	_, err = l.svcCtx.Rdb.Hincrby(fmt.Sprintf("inventory:%d", in.ProductId), "total", int(in.Quantity))

	if err != nil {
		l.Logger.Errorw("update inventory total failed", logx.Field("product_id", in.ProductId), logx.Field("err", err))
		return nil, err
	}

	//执行sql增加库存
	Total, err := l.svcCtx.InventoryModel.ReturnInventory(l.ctx, int32(in.ProductId), int32(in.Quantity))

	switch {
	case errors.Is(err, sqlx.ErrNotFound):
		l.Logger.Infow("product not in inventory", logx.Field("product_id", in.ProductId))
		res.StatusCode = code.ProductNotFoundInventory
		res.StatusMsg = code.ProductNotFoundInventoryMsg
		return res, nil

	case errors.Is(err, biz.InventoryDecreaseFailedErr):
		l.Logger.Errorw("inventory return failed", logx.Field("product_id", in.ProductId))
		return nil, err
	}

	res.Inventory = int64(Total)
	return res, nil

}
