package logic

import (
	"context"
	"fmt"
	"strconv"

	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecreasePreInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecreasePreInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecreasePreInventoryLogic {
	return &DecreasePreInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DecreaseInventory 预扣减库存，此时并非真实扣除库存，而是在缓存进行--操作
func (l *DecreasePreInventoryLogic) DecreasePreInventory(in *inventory.InventoryReq) (*inventory.InventoryResp, error) {

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

	totalInt, err := strconv.Atoi(total)

	if err != nil {
		l.Logger.Errorw("convert inventory total to int failed", logx.Field("product_id", in.ProductId), logx.Field("err", err))
		return nil, err
	}

	// 判断库存是否足够
	if in.Quantity > int32(totalInt) {
		l.Logger.Infow("product inventory not enough", logx.Field("product_id", in.ProductId), logx.Field("total", totalInt))
		res.StatusCode = code.InventoryNotEnough
		res.StatusMsg = code.InventoryNotEnoughMsg
		return res, nil
	}

	// 使用原子操作减少库存
	newTotal, err := l.svcCtx.Rdb.Hincrby(fmt.Sprintf("inventory:%d", in.ProductId), "total", -int(in.Quantity))

	if err != nil {
		l.Logger.Errorw("update inventory total failed", logx.Field("product_id", in.ProductId), logx.Field("err", err))
		return nil, err
	}

	if newTotal < 0 {
		// 检查库存是否被减到负数，如果是，应该恢复库存并返回错误
		_, err = l.svcCtx.Rdb.Hincrby(fmt.Sprintf("inventory:%d", in.ProductId), "total", int(in.Quantity))
		if err != nil {
			l.Logger.Errorw("rollback inventory total failed", logx.Field("product_id", in.ProductId), logx.Field("err", err))
			return nil, err
		}
		l.Logger.Infow("product inventory not enough", logx.Field("product_id", in.ProductId), logx.Field("total", newTotal+int(in.Quantity)))
		res.StatusCode = code.InventoryNotEnough
		res.StatusMsg = code.InventoryNotEnoughMsg
		return res, nil
	}

	res.Inventory = int64(newTotal)
	return res, nil
}
