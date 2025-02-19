package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	"jijizhazha1024/go-mall/services/inventory/inventory"
	"sync"

	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllProductLogic {
	return &GetAllProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页得到全部商品
func (l *GetAllProductLogic) GetAllProduct(in *product.GetAllProductsReq) (*product.GetAllProductsResp, error) {

	// 并发查询数据
	var wg sync.WaitGroup
	var products []*product2.Products
	var total int64
	var queryErr error
	productModel := product2.NewProductsModel(l.svcCtx.Mysql)
	wg.Add(2)

	// 查询商品列表
	go func() {
		defer wg.Done()
		offset := (in.Page - 1) * in.PageSize
		products, queryErr = productModel.FindPage(l.ctx, int(offset), int(in.PageSize))
	}()

	// 查询总数
	go func() {
		defer wg.Done()
		total, queryErr = productModel.Count(l.ctx)
	}()
	wg.Wait()

	// 统一错误处理
	if queryErr != nil {
		if errors.Is(queryErr, sqlx.ErrNotFound) {
			// 也可以记录info日志
			// 返回空，可能是由于用户的过滤条件导致没有匹配到数据

			return &product.GetAllProductsResp{}, nil
		}
		l.Logger.Errorw("query products failed",
			logx.Field("err", queryErr),
			logx.Field("page", in.Page),
			logx.Field("pageSize", in.PageSize))
		return nil, queryErr
	}

	// TODO 这里可能需要使用 协程池优化
	// 预分配切片容量
	var wgStock sync.WaitGroup
	var wgCategories sync.WaitGroup
	result := make([]*product.Product, len(products))
	wgStock.Add(len(products))
	wgCategories.Add(len(products))
	for i, p := range products {
		result[i] = &product.Product{
			Id:          uint32(p.Id),
			Name:        p.Name,
			Description: p.Description.String,
			Picture:     p.Picture.String,
			Price:       p.Price,
		}
		go func(index int, productId int64) {
			defer wgStock.Done()
			// 调用库存服务
			inventoryResp, err := l.svcCtx.InventoryRpc.GetInventory(l.ctx, &inventory.GetInventoryReq{
				ProductId: int32(productId),
			})
			if err != nil {
				l.Logger.Errorw("call rpc InventoryRpc.GetInventory failed", logx.Field("err", err), logx.Field("product_id", productId))
				return // 返回默认值或特殊标记
			}
			// 安全更新库存信息
			result[index].Stock = inventoryResp.Inventory
			result[index].Sold = inventoryResp.SoldCount
		}(i, p.Id) // 注意这里要显式传递参数

		go func(index int, productId int64) {
			defer wgCategories.Done()
			categories, err := l.svcCtx.CategoriesModel.FindCategoryNameByProductID(l.ctx, productId)
			if err != nil {
				l.Logger.Errorw("Failed to find product_category from database",
					logx.Field("err", err),
					logx.Field("product_id", productId))
				// 因为查询不完整，所以不需要写入缓存了，直接返回
				return
			}
			result[index].Categories = categories
		}(i, p.Id)
	}
	// 等待所有goroutine完成
	wgStock.Wait()
	wgCategories.Wait()
	return &product.GetAllProductsResp{
		Products: result,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}
