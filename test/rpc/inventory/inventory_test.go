package inventory

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/inventory/inventory"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var invClient inventory.InventoryClient
var once2 sync.Once

func setupInventoryClient(t *testing.T) {
	once2.Do(func() {
		conn, err := grpc.NewClient(
			fmt.Sprintf("127.0.0.1:%d", biz.InventoryRpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			t.Fatalf("连接库存服务失败: %v", err)
		}
		invClient = inventory.NewInventoryClient(conn)
	})
}

func TestInventoryService(t *testing.T) {

	t.Run("预扣库存全流程", func(t *testing.T) {
		setupInventoryClient(t)
		ctx := context.Background()
		testProductID := int32(1) // 测试用商品ID
		testProductID2 := int32(2)
		// testPreOrderID := "PRE_ORDER_2" // 测试用预订单ID

		// testuserID := int32(2) // 测试用用户ID
		//	初始化库存1
		_, err := invClient.UpdateInventory(ctx, &inventory.UpdateInventoryReq{
			Items: []*inventory.UpdateInventoryReq_Items{
				{ProductId: testProductID, Quantity: 100},
				{ProductId: testProductID2, Quantity: 200},
			},
		})
		assert.NoError(t, err)

		// // 预扣库存
		// preDecResp, err := invClient.DecreasePreInventory(ctx, &inventory.InventoryReq{
		// 	Items: []*inventory.InventoryReq_Items{
		// 		{ProductId: testProductID, Quantity: 30},
		// 	},
		// 	PreOrderId: testPreOrderID,
		// 	UserId:     testuserID,
		// })
		// assert.NoError(t, err)
		// fmt.Println("err---------------------------------------", err)
		// fmt.Println("preDecResp-", preDecResp)

		// //真实扣减
		// realDecResp, err := invClient.DecreaseInventory(ctx, &inventory.InventoryReq{
		// 	Items: []*inventory.InventoryReq_Items{
		// 		{ProductId: testProductID, Quantity: 30},
		// 	},
		// 	PreOrderId: testPreOrderID,
		// 	UserId:     testuserID,
		// })
		// assert.NoError(t, err)
		// fmt.Println("err-------------------", err)
		// fmt.Println("realDecResp-", realDecResp)
		// //第二次真实扣减
		// realDecResp, err = invClient.DecreaseInventory(ctx, &inventory.InventoryReq{
		// 	Items: []*inventory.InventoryReq_Items{
		// 		{ProductId: testProductID, Quantity: 30},
		// 	},
		// 	PreOrderId: testPreOrderID,
		// 	UserId:     testuserID,
		// })
		// assert.NoError(t, err)
		// fmt.Println("err-------------------", err)
		// fmt.Println("realDecResp-", realDecResp)

		// //归还缓存库存
		// retResp, err := invClient.ReturnPreInventory(ctx, &inventory.InventoryReq{
		// 	Items: []*inventory.InventoryReq_Items{
		// 		{ProductId: testProductID, Quantity: 30},
		// 	},
		// 	PreOrderId: testPreOrderID,
		// 	UserId:     testuserID,
		// })
		// assert.NoError(t, err)
		// fmt.Println("err-", err)
		// fmt.Println("retResp-", retResp)
		// //归还真实库存
		// retResp, err = invClient.ReturnInventory(ctx, &inventory.InventoryReq{
		// 	Items: []*inventory.InventoryReq_Items{
		// 		{ProductId: testProductID, Quantity: 30},
		// 	},
		// 	PreOrderId: testPreOrderID,
		// 	UserId:     testuserID,
		// })
		// assert.NoError(t, err)
		// fmt.Println("err-", err)
		// fmt.Println("retResp-", retResp)
		// //第二次归还真实库存
		// retResp, err = invClient.ReturnInventory(ctx, &inventory.InventoryReq{
		// 	Items: []*inventory.InventoryReq_Items{
		// 		{ProductId: testProductID, Quantity: 30},
		// 	},
		// 	PreOrderId: testPreOrderID,
		// 	UserId:     testuserID,
		// })
		// assert.NoError(t, err)
		// fmt.Println("err-", err)
		// fmt.Println("retResp-", retResp)

		// 验证最终库存
		getResp, err := invClient.GetInventory(ctx, &inventory.GetInventoryReq{
			ProductId: testProductID,
		})
		assert.NoError(t, err)
		fmt.Println("err-", err)
		fmt.Println("getResp-", getResp)

		getResp, err = invClient.GetInventory(ctx, &inventory.GetInventoryReq{
			ProductId: testProductID2,
		})
		assert.NoError(t, err)
		fmt.Println("err-", err)
		fmt.Println("getResp-", getResp)
	})

	// 	t.Run("异常场景处理", func(t *testing.T) {
	// 		// 无效商品ID
	// 		_, err := invClient.GetInventory(ctx, &inventory.GetInventoryReq{
	// 			ProductId: 9999,
	// 		})
	// 		assert.Error(t, err)
	// 		assert.Contains(t, err.Error(), "not found")

	// 		// 负数库存更新
	// 		_, err = invClient.UpdateInventory(ctx, &inventory.InventoryReq{
	// 			Items: []*inventory.InventoryReq_Items{
	// 				{ProductId: testProductID, Quantity: -100},
	// 			},
	// 		})
	// 		assert.Error(t, err)
	// 		assert.Contains(t, err.Error(), "invalid quantity")

	// 		// 预扣超量
	// 		_, err = invClient.DecreasePreInventory(ctx, &inventory.InventoryReq{
	// 			Items: []*inventory.InventoryReq_Items{
	// 				{ProductId: testProductID, Quantity: 1000},
	// 			},
	// 		})
	// 		assert.Error(t, err)
	// 		assert.Contains(t, err.Error(), "insufficient inventory")
	// 	})

	// 	t.Run("分布式锁竞争测试", func(t *testing.T) {
	// 		const (
	// 			concurrentRequests = 50
	// 			initialStock       = 500
	// 			deductPerRequest   = 10
	// 		)

	// 		// 初始化库存
	// 		_, err := invClient.UpdateInventory(ctx, &inventory.InventoryReq{
	// 			Items: []*inventory.InventoryReq_Items{
	// 				{ProductId: 5001, Quantity: int32(initialStock)},
	// 			},
	// 		})
	// 		assert.NoError(t, err)

	// 		var wg sync.WaitGroup
	// 		wg.Add(concurrentRequests)

	// 		// 并发扣减
	// 		for i := 0; i < concurrentRequests; i++ {
	// 			go func() {
	// 				defer wg.Done()
	// 				_, err := invClient.DecreaseInventory(ctx, &inventory.InventoryReq{
	// 					Items: []*inventory.InventoryReq_Items{
	// 						{ProductId: 5001, Quantity: deductPerRequest},
	// 					},
	// 				})
	// 				assert.NoError(t, err)
	// 			}()
	// 		}
	// 		wg.Wait()

	// 		// 验证最终库存
	// 		resp, err := invClient.GetInventory(ctx, &inventory.GetInventoryReq{
	// 			ProductId: 5001,
	// 		})
	// 		assert.NoError(t, err)
	// 		expected := initialStock - concurrentRequests*deductPerRequest
	// 		assert.Equal(t, int64(expected), resp.Inventory)
	// 	})

	// 	t.Run("库存归还测试", func(t *testing.T) {
	// 		// 初始化
	// 		testProduct := int32(6001)
	// 		_, err := invClient.UpdateInventory(ctx, &inventory.InventoryReq{
	// 			Items: []*inventory.InventoryReq_Items{
	// 				{ProductId: testProduct, Quantity: 100},
	// 			},
	// 		})
	// 		assert.NoError(t, err)

	// 		// 扣减库存
	// 		_, err = invClient.DecreaseInventory(ctx, &inventory.InventoryReq{
	// 			Items: []*inventory.InventoryReq_Items{
	// 				{ProductId: testProduct, Quantity: 50},
	// 			},
	// 		})
	// 		assert.NoError(t, err)

	// 		// 部分归还
	// 		retResp, err := invClient.ReturnInventory(ctx, &inventory.InventoryReq{
	// 			Items: []*inventory.InventoryReq_Items{
	// 				{ProductId: testProduct, Quantity: 20},
	// 			},
	// 		})
	// 		assert.NoError(t, err)
	// 		assert.Equal(t, code.Success, retResp.StatusCode)

	// 		// 验证结果
	// 		getResp, err := invClient.GetInventory(ctx, &inventory.GetInventoryReq{
	// 			ProductId: testProduct,
	// 		})
	// 		assert.NoError(t, err)
	// 		assert.Equal(t, int64(70), getResp.Inventory)
	// 	})
	// }

	// // 单独测试接口
	// func TestGetInventory_NotFound(t *testing.T) {
	// 	setupInventoryClient(t)

	// 	resp, err := invClient.GetInventory(context.Background(), &inventory.GetInventoryReq{
	// 		ProductId: 9999, // 不存在的商品
	// 	})
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, code.ProductNotFoundInventory, resp.StatusCode)
	// }

	// // func TestReturnInventory_InvalidParam(t *testing.T) {
	// // 	setupInventoryClient(t)

	// // 	// 测试负数参数
	// // 	_, err := invClient.ReturnPreInventory(context.Background(), &inventory.InventoryReq{
	// // 		ProductId: 1001,
	// // 		Quantity:  -10,
	// // 	})
	// // 	assert.Error(t, err)
	// // 	assert.Contains(t, err.Error(), "invalid inventory")
}
func TestInventoryService_HighConcurrency(t *testing.T) {
	setupInventoryClient(t)
	ctx := context.Background()

	// 测试商品参数
	testProductID := int32(9999)
	initialStock := int32(1000)    // 初始库存
	concurrentUsers := 100         // 并发用户数
	deductPerUser := int32(10)     // 每个用户扣减量
	expectedFinalStock := int32(0) // 预期最终库存

	// 初始化库存
	_, err := invClient.UpdateInventory(ctx, &inventory.UpdateInventoryReq{
		Items: []*inventory.UpdateInventoryReq_Items{
			{ProductId: testProductID, Quantity: initialStock},
		},
	})
	assert.NoError(t, err, "初始化库存失败")

	var wg sync.WaitGroup
	wg.Add(concurrentUsers)

	// 模拟并发请求
	for i := 0; i < concurrentUsers; i++ {
		go func(userID int32) {
			defer wg.Done()

			// 生成唯一预订单ID
			preOrderID := fmt.Sprintf("PRE_ORDER_%d_%d", userID, time.Now().UnixNano())

			// 步骤1: 预扣库存
			_, preErr := invClient.DecreasePreInventory(ctx, &inventory.InventoryReq{
				Items: []*inventory.InventoryReq_Items{
					{ProductId: testProductID, Quantity: deductPerUser},
				},
				PreOrderId: preOrderID,
				UserId:     userID,
			})
			if preErr != nil {
				t.Logf("用户 %d 预扣失败: %v", userID, preErr)
				return
			}

			// 步骤2: 真实扣减
			_, decErr := invClient.DecreaseInventory(ctx, &inventory.InventoryReq{
				Items: []*inventory.InventoryReq_Items{
					{ProductId: testProductID, Quantity: deductPerUser},
				},
				PreOrderId: preOrderID,
				UserId:     userID,
			})
			if decErr != nil {
				t.Logf("用户 %d 扣减失败: %v", userID, decErr)
				// 失败后归还预扣库存
				_, retErr := invClient.ReturnPreInventory(ctx, &inventory.InventoryReq{
					Items: []*inventory.InventoryReq_Items{
						{ProductId: testProductID, Quantity: deductPerUser},
					},
					PreOrderId: preOrderID,
					UserId:     userID,
				})
				assert.NoError(t, retErr, "预扣库存归还失败")
			}
		}(int32(i + 400)) // 生成唯一用户ID
	}

	// 等待所有请求完成
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// 设置超时时间
	select {
	case <-done:
		t.Log("所有并发请求完成")
	case <-time.After(50 * time.Second):
		t.Fatal("测试超时：50秒未完成")
	}

	// 验证最终库存
	finalStockResp, err := invClient.GetInventory(ctx, &inventory.GetInventoryReq{
		ProductId: testProductID,
	})
	assert.NoError(t, err, "获取最终库存失败")

	t.Logf("预期库存: %d, 实际库存: %d",
		expectedFinalStock, finalStockResp.Inventory)

	assert.Equal(t, expectedFinalStock, finalStockResp.Inventory,
		"最终库存不一致，可能存在并发问题")

	// 验证总扣减次数（可选）
	// 可通过日志系统或监控指标验证实际扣减次数是否为 concurrentUsers
}
