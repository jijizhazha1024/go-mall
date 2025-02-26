package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"
	"sync"

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

// GetAllProduct 分页得到全部商品
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

	// 预分配切片容量
	result := make([]*product.Product, len(products))
	populateProductDetails(l.ctx, l.svcCtx, products, result)
	return &product.GetAllProductsResp{
		Products: result,
		Total:    total,
		Page:     in.Page,
		PageSize: in.PageSize,
	}, nil
}
