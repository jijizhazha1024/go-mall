package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"jijizhazha1024/go-mall/services/checkout/internal/svc"
	"jijizhazha1024/go-mall/services/coupons/coupons"
	"jijizhazha1024/go-mall/services/inventory/inventory"
	"jijizhazha1024/go-mall/services/product/product"
	"strings"
	"time"
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

// PrepareCheckout 处理预结算
func (l *PrepareCheckoutLogic) PrepareCheckout(in *checkout.CheckoutReq) (*checkout.CheckoutResp, error) {
	// 1. 生成 pre_order_id
	preOrderId, err := generatePreOrderID()
	if err != nil {
		l.Logger.Errorw("生成 preOrderId 失败", logx.Field("err", err), logx.Field("user_id", in.UserId))
		return nil, errors.New("生成订单ID失败")
	}

	// 2. 使用 Redis 锁来保证幂等性
	cacheKey := fmt.Sprintf("checkout:preorder:%d", in.UserId)
	luaScript := `
		if redis.call("EXISTS", KEYS[1]) == 1 then
			return 0
		else
			redis.call("SETEX", KEYS[1], ARGV[1], ARGV[2])
			return 1
		end
	`
	result, err := l.svcCtx.RedisClient.Eval(luaScript, []string{cacheKey}, preOrderId, 300)
	if err != nil {
		l.Logger.Errorw("Redis Lua 执行失败", logx.Field("err", err), logx.Field("user_id", in.UserId))
		return nil, errors.New("系统错误")
	}
	if result == int64(0) {
		l.Logger.Infof("用户 %d 的预订单 %s 已存在，跳过重复结算", in.UserId, preOrderId)
		return &checkout.CheckoutResp{
			StatusCode: 200,
			StatusMsg:  "预结算已处理",
			PreOrderId: preOrderId,
		}, nil
	}

	// 3. 检查是否有商品信息
	if len(in.OrderItems) == 0 {
		// 释放 Redis 锁
		if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
			l.Logger.Errorw("删除 Redis 锁失败", logx.Field("err", err), logx.Field("user_id", in.UserId))
		}
		return &checkout.CheckoutResp{
			StatusCode: 400,
			StatusMsg:  "订单商品不能为空",
		}, nil
	}

	// 4. 调用库存预扣接口
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
		l.Logger.Errorw("库存预扣失败，执行同步库存回滚", logx.Field("err", err), logx.Field("user_id", in.UserId), logx.Field("pre_order_id", preOrderId))

		// 释放 Redis 锁
		if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
			l.Logger.Errorw("删除 Redis 锁失败", logx.Field("err", err), logx.Field("user_id", in.UserId))
		}

		// 同步回滚库存
		_, errRollback := l.svcCtx.InventoryRpc.ReturnPreInventory(l.ctx, &inventory.InventoryReq{
			Items:      inventoryItems,
			PreOrderId: preOrderId,
			UserId:     int32(in.UserId),
		})
		if errRollback != nil {
			l.Logger.Errorw("库存回滚失败", logx.Field("err", errRollback), logx.Field("user_id", in.UserId), logx.Field("pre_order_id", preOrderId))
		}

		return nil, errors.New("库存不足")
	}

	// 5. 异步处理结算信息
	go func() {
		err := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
			var totalOriginalAmount int64
			var finalAmount int64
			var orderItems []*coupons.Items
			var availableCoupons []string

			// 1. 删除购物车中的商品
			for _, item := range in.OrderItems {
				_, err := session.Exec("DELETE FROM carts WHERE user_id = ? AND product_id = ?", in.UserId, item.ProductId)
				if err != nil {
					return err
				}
			}

			// 2. 获取商品信息，计算原始总金额并插入 checkout_items
			for _, item := range in.OrderItems {
				productResp, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.GetProductReq{
					Id: uint32(item.ProductId),
				})
				if err != nil || productResp.Product == nil {
					l.Logger.Errorw("获取商品详情失败", logx.Field("err", err), logx.Field("product_id", item.ProductId))
					return errors.New("获取商品信息失败")
				}

				snapshotData := map[string]interface{}{"name": productResp.Product.Name, "specs": productResp.Product.Description}
				snapshotJSON, _ := json.Marshal(snapshotData)

				// 插入 checkout_items 表
				_, err = session.Exec(
					"INSERT INTO checkout_items (pre_order_id, product_id, quantity, price, snapshot, created_at) VALUES (?, ?, ?, ?, ?, NOW())",
					preOrderId, item.ProductId, item.Quantity, productResp.Product.Price, string(snapshotJSON),
				)
				if err != nil {
					return err
				}

				// 累加商品原始总金额
				totalOriginalAmount += productResp.Product.Price * int64(item.Quantity)
				orderItems = append(orderItems, &coupons.Items{
					ProductId: item.ProductId,
					Quantity:  item.Quantity,
				})
			}

			userCouponsResp, err := l.svcCtx.CouponsRpc.ListUserCoupons(l.ctx, &coupons.ListUserCouponsReq{
				UserId:     int32(in.UserId),
				Pagination: &coupons.PaginationReq{Page: 1, Size: 50},
			})
			if err != nil {
				l.Logger.Errorw("查询用户优惠券失败", logx.Field("err", err))
				return errors.New("查询用户优惠券失败")
			}

			finalAmount = totalOriginalAmount
			for _, coupon := range userCouponsResp.UserCoupons {
				couponResp, err := l.svcCtx.CouponsRpc.CalculateCoupon(l.ctx, &coupons.CalculateCouponReq{
					UserId:   int32(in.UserId),
					CouponId: coupon.CouponId,
					Items:    orderItems,
				})
				if err != nil {
					l.Logger.Errorw("计算优惠失败", logx.Field("err", err), logx.Field("coupon_id", coupon.Id))
					continue
				}

				// 如果优惠券可用，则应用折扣
				if couponResp.IsUsable {
					finalAmount -= couponResp.DiscountAmount
					availableCoupons = append(availableCoupons, coupon.CouponId)
				}
			}

			// 确保最终金额不小于 0
			if finalAmount < 0 {
				finalAmount = 0
			}

			_, err = session.Exec(
				"INSERT INTO checkouts (pre_order_id, user_id, coupon_id, original_amount, final_amount, status, expire_time, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())",
				preOrderId, in.UserId, strings.Join(availableCoupons, ","), totalOriginalAmount, finalAmount, checkout.CheckoutStatus_RESERVING, time.Now().Add(10*time.Minute).Unix(),
			)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			l.Logger.Errorw("处理结算信息失败", logx.Field("err", err))
		}
	}()

	// 释放 Redis 锁
	if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
		l.Logger.Errorw("删除 Redis 锁失败", logx.Field("err", err), logx.Field("user_id", in.UserId))
	}

	// 6. 返回预结算信息
	return &checkout.CheckoutResp{
		StatusCode: 200,
		StatusMsg:  "预结算成功",
		PreOrderId: preOrderId,
		ExpireTime: time.Now().Add(10 * time.Minute).Unix(),
		PayMethod:  []int64{1, 2},
	}, nil
}
