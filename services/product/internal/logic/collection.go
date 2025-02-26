package logic

import (
	"context"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/zeromicro/go-zero/core/logx"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	"jijizhazha1024/go-mall/services/inventory/inventory"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"
	"sync"
)

var pool = gopool.NewPool("product-details-pool", 100, gopool.NewConfig()) // 根据实际情况调整参数
// 通用并发处理库存和分类信息
func populateProductDetails(ctx context.Context, svcCtx *svc.ServiceContext, products []*product2.Products, result []*product.Product) {
	var wg sync.WaitGroup
	wg.Add(2 * len(products))
	for index, p := range products {
		pool.Go(func() {
			handleInventory(ctx, svcCtx, result, index, p.Id)
		})
		pool.Go(func() {
			handleCategories(ctx, svcCtx, result, index, p.Id)
		})
	}
	wg.Wait()
}

// 库存处理逻辑
func handleInventory(ctx context.Context, svcCtx *svc.ServiceContext, result []*product.Product, index int, productId int64) {
	inventoryResp, err := svcCtx.InventoryRpc.GetInventory(ctx, &inventory.GetInventoryReq{
		ProductId: int32(productId),
	})
	if err != nil {
		logx.WithContext(ctx).Errorw("call InventoryRpc failed", logx.Field("err", err), logx.Field("product_id", productId))
		return
	}
	result[index].Stock = inventoryResp.Inventory
	result[index].Sold = inventoryResp.SoldCount
}

// 分类处理逻辑
func handleCategories(ctx context.Context, svcCtx *svc.ServiceContext, result []*product.Product, index int, productId int64) {
	categories, err := svcCtx.CategoriesModel.FindCategoryNameByProductID(ctx, productId)
	if err != nil {
		logx.WithContext(ctx).Errorw("query categories failed", logx.Field("err", err), logx.Field("product_id", productId))
		return
	}
	result[index].Categories = categories
}
