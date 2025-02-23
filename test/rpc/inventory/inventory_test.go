package inventory

// import (
// 	"context"
// 	"fmt"
// 	"jijizhazha1024/go-mall/common/consts/biz"
// 	"jijizhazha1024/go-mall/common/consts/code"
// 	"jijizhazha1024/go-mall/services/inventory/inventory"
// 	"sync"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// var invClient inventory.InventoryClient
// var once2 sync.Once

// func setupInventoryClient(t *testing.T) {
// 	once2.Do(func() {
// 		conn, err := grpc.NewClient(
// 			fmt.Sprintf("127.0.0.1:%d", biz.InventoryRpcPort),
// 			grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		)
// 		if err != nil {
// 			t.Fatalf("连接库存服务失败: %v", err)
// 		}
// 		invClient = inventory.NewInventoryClient(conn)
// 	})
// }

// func TestInventoryService(t *testing.T) {
// 	setupInventoryClient(t)
// 	ctx := context.Background()
// 	testProductID := int32(1000) // 测试用商品ID

// 	t.Run("基本库存操作流程", func(t *testing.T) {
// 		// 1. 初始化库存
// 		initResp, err := invClient.UpdateInventory(ctx, &inventory.InventoryReq{
// 			ProductId: testProductID,
// 			Quantity:  100,
// 		})
// 		assert.NoError(t, err)
// 		fmt.Println("err-----------------------", err)
// 		fmt.Println("resp-----------------------", initResp)
// 		assert.Equal(t, int64(100), initResp.Inventory)

// 		// 2. 查询库存
// 		getResp, err := invClient.GetInventory(ctx, &inventory.GetInventoryReq{
// 			ProductId: testProductID,
// 		})
// 		assert.NoError(t, err)

// 		assert.Equal(t, int64(100), getResp.Inventory)

// 		// 3. 扣减库存
// 		decResp, err := invClient.DecreaseInventory(ctx, &inventory.InventoryReq{
// 			ProductId: testProductID,
// 			Quantity:  20,
// 		})
// 		assert.NoError(t, err)

// 		assert.Equal(t, int64(80), decResp.Inventory)

// 		// 4. 归还库存
// 		retResp, err := invClient.ReturnPreInventory(ctx, &inventory.InventoryReq{
// 			ProductId: testProductID,
// 			Quantity:  10,
// 		})
// 		assert.NoError(t, err)

// 		assert.Equal(t, int64(90), retResp.Inventory)
// 	})

// 	t.Run("边界值测试", func(t *testing.T) {
// 		// 扣减超过库存量
// 		res, err := invClient.DecreaseInventory(ctx, &inventory.InventoryReq{
// 			ProductId: testProductID,
// 			Quantity:  1000,
// 		})
// 		assert.NoError(t, err)
// 		assert.Equal(t, res.StatusCode, code.InventoryNotEnough)
// 	})

// 	t.Run("并发扣减测试", func(t *testing.T) {
// 		const parallel = 10
// 		var wg sync.WaitGroup
// 		wg.Add(parallel)

// 		// 初始化库存
// 		_, _ = invClient.UpdateInventory(ctx, &inventory.InventoryReq{
// 			ProductId: testProductID,
// 			Quantity:  100,
// 		})

// 		for i := 0; i < parallel; i++ {
// 			go func() {
// 				defer wg.Done()
// 				_, err := invClient.DecreaseInventory(ctx, &inventory.InventoryReq{
// 					ProductId: testProductID,
// 					Quantity:  10,
// 				})
// 				assert.NoError(t, err)
// 			}()
// 		}
// 		wg.Wait()

// 		// 验证最终库存
// 		resp, err := invClient.GetInventory(ctx, &inventory.GetInventoryReq{
// 			ProductId: testProductID,
// 		})
// 		assert.NoError(t, err)
// 		assert.Equal(t, int64(0), resp.Inventory)
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

// func TestReturnInventory_InvalidParam(t *testing.T) {
// 	setupInventoryClient(t)

// 	// 测试负数参数
// 	_, err := invClient.ReturnPreInventory(context.Background(), &inventory.InventoryReq{
// 		ProductId: 1001,
// 		Quantity:  -10,
// 	})
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "invalid inventory")
// }
