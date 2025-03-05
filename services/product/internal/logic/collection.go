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
	"time"
)

var pool = gopool.NewPool("product-details-pool", 100, gopool.NewConfig()) // 根据实际情况调整参数
// 通用并发处理库存和分类信息
func populateProductDetails(ctx context.Context, svcCtx *svc.ServiceContext, products []*product2.Products) (result []*product.Product) {
	result = make([]*product.Product, len(products))
	var wg sync.WaitGroup
	wg.Add(len(products)) // 每个产品一个任务

	for i := range products {
		index := i // 创建局部变量避免闭包问题
		p := products[index]
		result[index] = &product.Product{
			Id:          uint32(p.Id),
			Name:        p.Name,
			Description: p.Description.String,
			Picture:     p.Picture.String,
			Price:       p.Price,
			CratedAt:    p.CreatedAt.Format(time.DateTime),
			UpdatedAt:   p.UpdatedAt.Format(time.DateTime),
		} // 初始化结果
		pool.CtxGo(ctx, func() {
			defer wg.Done()
			var innerWg sync.WaitGroup
			innerWg.Add(2)
			// 处理库存
			go func() {
				defer innerWg.Done()
				handleInventory(ctx, svcCtx, result, index, p.Id)
			}()

			// 处理分类
			go func() {
				defer innerWg.Done()
				handleCategories(ctx, svcCtx, result, index, p.Id)
			}()

			innerWg.Wait()
		})
	}

	wg.Wait()
	return result
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
