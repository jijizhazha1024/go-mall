package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/services/carts/carts"
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

// PrepareCheckout 预结算 (幂等删除购物车)
func (l *PrepareCheckoutLogic) PrepareCheckout(in *checkout.CheckoutReq) (*checkout.CheckoutResp, error) {
	// 1. 生成 pre_order_id
	preOrderId, err := generatePreOrderID()
	if err != nil {
		logx.Errorf("生成 preOrderId 失败: %v", err)
		return nil, errors.New("生成订单ID失败")
	}

	// 2. 幂等性检查 (Redis 幂等锁)
	cacheKey := fmt.Sprintf("checkout:preorder:%d", in.UserId)
	exist, err := l.svcCtx.RedisClient.Exists(cacheKey)
	if err != nil {
		logx.Errorf("Redis 查询失败: %v", err)
		return nil, errors.New("系统错误")
	}

	if exist { // 如果该用户已有预结算订单，直接返回，防止重复删除购物车
		logx.Infof("用户 %d 的预订单 %s 已存在，跳过重复结算", in.UserId, preOrderId)
		return &checkout.CheckoutResp{
			StatusCode: 200,
			StatusMsg:  "预结算已处理",
			PreOrderId: preOrderId,
		}, nil
	}

	// 3. 设置 Redis 幂等锁 (有效期 5 分钟)
	_, err = l.svcCtx.RedisClient.SetnxEx(cacheKey, preOrderId, 300) // 5 分钟过期
	if err != nil {
		logx.Errorf("Redis 设置失败: %v", err)
		return nil, errors.New("系统错误")
	}

	// 4. 查询购物车
	cartItems, err := l.svcCtx.CartsModel.FindByUserID(l.ctx, int64(in.UserId))
	if err != nil {
		logx.Errorf("获取购物车失败: user_id=%d, err=%v", in.UserId, err)
		return nil, errors.New("获取购物车失败")
	}

	if len(cartItems) == 0 {
		return &checkout.CheckoutResp{
			StatusCode: 400,
			StatusMsg:  "购物车为空",
		}, nil
	}

	// 5. 调用库存预扣服务
	inventoryItems := make([]*inventory.InventoryReq_Items, 0)
	for _, item := range cartItems {
		inventoryItems = append(inventoryItems, &inventory.InventoryReq_Items{
			ProductId: int32(item.ProductId.Int64),
			Quantity:  int32(item.Quantity.Int64),
		})
	}

	_, err = l.svcCtx.InventoryRpc.DecreasePreInventory(l.ctx, &inventory.InventoryReq{
		Items:      inventoryItems,
		PreOrderId: preOrderId,
		UserId:     int32(in.UserId),
	})

	if err != nil {
		logx.Errorf("库存预扣失败: user_id=%d, pre_order_id=%s, err=%v", in.UserId, preOrderId, err)
		return nil, errors.New("库存不足")
	}

	// 6. 存储结算信息到 Redis
	_, err = l.svcCtx.RedisClient.SetnxEx(fmt.Sprintf("checkout:order:%s", preOrderId), preOrderId, 300)
	if err != nil {
		logx.Errorf("Redis 存储结算信息失败: %v", err)
		return nil, errors.New("系统错误")
	}

	// 7. 启动异步任务处理后续操作（删除购物车）
	go func() {
		// 删除购物车商品
		for _, item := range cartItems {
			_, err := l.svcCtx.CartsRpc.DeleteCartItem(l.ctx, &carts.CartItemRequest{
				UserId:    int32(in.UserId),
				ProductId: int32(item.ProductId.Int64),
			})

			if err != nil {
				logx.Errorf("删除购物车失败: user_id=%d, product_id=%d, err=%v", in.UserId, item.ProductId, err)
			} else {
				logx.Infof("成功删除购物车商品: user_id=%d, product_id=%d", in.UserId, item.ProductId)
			}
		}
	}()

	// 8. 返回预结算信息
	return &checkout.CheckoutResp{
		StatusCode: 200,
		StatusMsg:  "预结算成功",
		PreOrderId: preOrderId,
	}, nil
}
