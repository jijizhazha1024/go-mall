package logic

import (
	"context"
	"errors"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

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

	//访问记录

	lockKey := fmt.Sprintf("%s:%d", biz.InventoryAccessKeyPrefix, in.ProductId)
	newcount, err := l.svcCtx.Rdb.Incr(lockKey)
	if err != nil {
		l.Logger.Infow("redis inventory incr failed", logx.Field("lock_key", lockKey))
		return nil, nil
	}
	// 设值过期时间，例如设置为4小时
	expireDuration := 4 * time.Hour
	expireInSeconds := int(expireDuration.Seconds())
	err = l.svcCtx.Rdb.Expire(lockKey, expireInSeconds)
	if err != nil {
		l.Logger.Infow("redis set expire failed", logx.Field("lock_key", lockKey))
		return nil, nil
	}

	if newcount > 500 {
		//将其加入缓存库存
		tostr := fmt.Sprintf("%d", inventoryResp.Total)
		err := l.svcCtx.Rdb.Set(fmt.Sprintf("%s:%d", biz.InventoryProductKey, in.ProductId), tostr)
		if err != nil {
			l.Logger.Infow("redis set failed", logx.Field("product_id", in.ProductId))
			return nil, nil
		}
		expireDuration := 1 * time.Hour
		expireInSeconds := int(expireDuration.Seconds())
		err = l.svcCtx.Rdb.Expire(fmt.Sprintf("%s:%d", biz.InventoryProductKey, in.ProductId), expireInSeconds)
		if err != nil {
			l.Logger.Infow("redis set expire failed", logx.Field("product_id", in.ProductId))
			return nil, nil
		}
	}

	res.Inventory = inventoryResp.Total
	res.SoldCount = inventoryResp.Sold
	return res, nil
}
