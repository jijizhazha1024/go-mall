package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"
)

type PrepareCheckoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPrepareCheckoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrepareCheckoutLogic {
	return &PrepareCheckoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func generatePreOrderID() (string, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// PrepareCheckout 预结算
func (l *PrepareCheckoutLogic) PrepareCheckout(in *checkout.CheckoutReq) (*checkout.CheckoutResp, error) {
	// 1. 生成 pre_order_id
	preOrderId, err := generatePreOrderID()
	if err != nil {
		l.Logger.Errorw("生成 preOrderId 失败",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId))
		return nil, errors.New("生成订单ID失败")
	}

	// 2. 使用 Lua 脚本设置 Redis 幂等锁
	cacheKey := fmt.Sprintf("checkout:preorder:%d", in.UserId)
	luaScript := `
		if redis.call("EXISTS", KEYS[1]) == 1 then
			return 0
		else
			redis.call("SETEX", KEYS[1], ARGV[1], ARGV[2])
			return 1
		end
	`
	result, err := l.svcCtx.RedisClient.Eval(luaScript, []string{cacheKey}, preOrderId, 300) // 5 分钟过期
	if err != nil {
		l.Logger.Errorw("Redis Lua 执行失败",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId))
		return nil, errors.New("系统错误")
	}

	if result == int64(0) { // key 已存在，说明有并发请求
		l.Logger.Infof("用户 %d 的预订单 %s 已存在，跳过重复结算", in.UserId, preOrderId)
		return &checkout.CheckoutResp{
			StatusCode: 200,
			StatusMsg:  "预结算已处理",
			PreOrderId: preOrderId,
		}, nil
	}

	// 3. 检查请求中是否传入商品信息
	if len(in.OrderItems) == 0 {
		return &checkout.CheckoutResp{
			StatusCode: 400,
			StatusMsg:  "订单商品不能为空",
		}, nil
	}

	// 4. 调用库存预扣服务
	inventoryItems := make([]*inventory.InventoryReq_Items, 0)
	for _, item := range in.OrderItems {
		inventoryItems = append(inventoryItems, &inventory.InventoryReq_Items{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	_, err = l.svcCtx.InventoryRpc.DecreasePreInventory(l.ctx, &inventory.InventoryReq{
		Items:      inventoryItems,
		PreOrderId: preOrderId,
		UserId:     int32(in.UserId),
	})

	if err != nil {
		l.Logger.Errorw("库存预扣失败，执行库存回滚",
			logx.Field("err", err),
			logx.Field("user_id", in.UserId),
			logx.Field("pre_order_id", preOrderId))

		// 库存预扣失败，调用退回预扣库存
		go l.svcCtx.InventoryRpc.ReturnPreInventory(l.ctx, &inventory.InventoryReq{
			Items:      inventoryItems,
			PreOrderId: preOrderId,
			UserId:     int32(in.UserId),
		})

		return nil, errors.New("库存不足")
	}

	// 5. 存储结算信息到 Redis
	_, err = l.svcCtx.RedisClient.SetnxEx(fmt.Sprintf("checkout:order:%s", preOrderId), preOrderId, 300)
	if err != nil {
		l.Logger.Errorw("Redis 存储结算信息失败，执行库存回滚",
			logx.Field("err", err),
			logx.Field("pre_order_id", preOrderId))

		// Redis 存储失败，回滚库存
		go l.svcCtx.InventoryRpc.ReturnPreInventory(l.ctx, &inventory.InventoryReq{
			Items:      inventoryItems,
			PreOrderId: preOrderId,
			UserId:     int32(in.UserId),
		})

		return nil, errors.New("系统错误")
	}

	// 6. 异步删除购物车
	go func() {
		err := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
			for _, item := range in.OrderItems {
				_, err := session.Exec("DELETE FROM carts WHERE user_id = ? AND product_id = ?", in.UserId, item.ProductId)
				if err != nil {
					l.Logger.Errorw("删除购物车商品失败",
						logx.Field("err", err),
						logx.Field("user_id", in.UserId),
						logx.Field("product_id", item.ProductId))
					return err
				}
			}
			return nil
		})

		if err != nil {
			l.Logger.Errorw("删除购物车事务失败",
				logx.Field("err", err),
				logx.Field("user_id", in.UserId))
		} else {
			l.Logger.Infow("删除购物车事务成功",
				logx.Field("user_id", in.UserId))
		}
	}()

	// 7. 返回预结算信息
	return &checkout.CheckoutResp{
		StatusCode: 200,
		StatusMsg:  "预结算成功",
		PreOrderId: preOrderId,
	}, nil
}
