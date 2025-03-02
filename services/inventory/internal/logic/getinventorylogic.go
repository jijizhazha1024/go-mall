package logic

import (
	"context"
	"errors"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"strconv"
	"time"

	"github.com/redis/rueidis/rueidislock"
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

	// 先从缓存中获取数据
	cacheKey := fmt.Sprintf("%s:%d", biz.InventoryProductKey, in.ProductId)
	total, err := l.svcCtx.Rdb.Get(cacheKey)
	if err == nil {
		inventoryResp := new(inventory.GetInventoryResp)
		inventoryResp.Inventory, _ = strconv.ParseInt(total, 10, 64)
		return inventoryResp, nil
	}

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

	accessKey := fmt.Sprintf(biz.InventoryAccessKeyPrefix, in.ProductId)
	newcount, err := l.svcCtx.Rdb.Incr(accessKey)
	if err != nil {
		l.Logger.Infow("redis inventory incr failed", logx.Field("lock_key", accessKey))
		return nil, nil
	}
	// 设值过期时间，例如设置为4小时
	expireDuration := 4 * time.Hour
	expireInSeconds := int(expireDuration.Seconds())
	err = l.svcCtx.Rdb.Expire(accessKey, expireInSeconds)
	if err != nil {
		l.Logger.Infow("redis set expire failed", logx.Field("lock_key", accessKey))
		return nil, nil
	}

	if newcount > 500 {
		// 创建分布式锁的key
		lockKey := fmt.Sprintf("%s:%d", biz.InventoryLockKey, in.ProductId)

		// 尝试获取非阻塞锁（立即返回）
		_, releaseLock, err := l.svcCtx.Locker.TryWithContext(l.ctx, lockKey)
		if err != nil {
			if errors.Is(err, rueidislock.ErrNotLocked) {
				// 锁已被占用，直接返回当前数据
				l.Logger.Infow("cache update in progress by another instance",
					logx.Field("product_id", in.ProductId))
				res.Inventory = inventoryResp.Total
				res.SoldCount = inventoryResp.Sold
				return res, nil
			}
			// 其他错误情况处理
			l.Logger.Errorw("failed to acquire distributed lock",
				logx.Field("error", err),
				logx.Field("product_id", in.ProductId))
			return nil, err
		}
		defer releaseLock()

		// 获取锁后再次检查访问量（双重检查）
		currentCount, err := l.svcCtx.Rdb.Get(accessKey)
		if err == nil && currentCount <= strconv.Itoa(500) {
			res.Inventory = inventoryResp.Total
			res.SoldCount = inventoryResp.Sold
			return res, nil
		}

		// 获取最新当前库存
		newinventoryResp, err := l.svcCtx.InventoryModel.FindOne(l.ctx, int64(in.ProductId))
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

		cacheKey := fmt.Sprintf("%s:%d", biz.InventoryProductKey, in.ProductId)
		total := strconv.Itoa(int(newinventoryResp.Total))
		if err := l.svcCtx.Rdb.Set(cacheKey, total); err != nil {
			l.Logger.Errorw("failed to update inventory cache",
				logx.Field("product_id", in.ProductId),
				logx.Field("error", err))
			return nil, err
		}

		// 设置缓存过期时间
		if err := l.svcCtx.Rdb.Expire(cacheKey, 3600); err != nil {
			l.Logger.Infow("failed to set cache expiration",
				logx.Field("product_id", in.ProductId),
				logx.Field("error", err))
		}

		// 重置访问计数器（原子操作）
		if err := l.svcCtx.Rdb.Set(accessKey, "0"); err != nil {
			l.Logger.Errorw("failed to reset access counter",
				logx.Field("product_id", in.ProductId),
				logx.Field("error", err))
		}

	}

	res.Inventory = inventoryResp.Total
	res.SoldCount = inventoryResp.Sold
	return res, nil
}
