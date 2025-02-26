package logic

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"
)

type QueryProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryProductLogic {
	return &QueryProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QueryProduct 根据条件查询商品
func (l *QueryProductLogic) QueryProduct(in *product.QueryProductReq) (*product.GetAllProductsResp, error) {
	// 1. 构建基础查询
	boolQuery := elastic.NewBoolQuery()
	if in.Name != "" {
		boolQuery.Must(elastic.NewMatchQuery("name", in.Name))
	}

	if in.Keyword != "" {
		boolQuery.Must(elastic.NewMultiMatchQuery(in.Keyword, "name", "description"))
	}

	if len(in.Category) > 0 {
		termsQuery := elastic.NewTermsQuery("category", stringSliceToInterface(in.Category)...)
		boolQuery.Must(termsQuery)
	}

	if in.Price != nil {
		rangeQuery := elastic.NewRangeQuery("price")
		if in.Price.Min > 0 {
			rangeQuery.Gte(in.Price.Min)
		}
		if in.Price.Max > 0 {
			rangeQuery.Lte(in.Price.Max)
		}
		boolQuery.Must(rangeQuery)
	}

	// 2. 分页处理
	from := 0
	size := biz.DefaultPageSize // 默认分页大小
	if in.Paginator != nil {
		// 处理分页大小
		if in.Paginator.PageSize > 0 {
			// 限制最大分页不超过30
			if in.Paginator.PageSize > biz.MaxPageSize {
				size = biz.MaxPageSize
			} else {
				size = int(in.Paginator.PageSize)
			}
		}
		// 计算起始位置（增加页码有效性校验）
		if in.Paginator.Page > 1 {
			from = (int(in.Paginator.Page) - 1) * size
		}
	}

	res := &product.GetAllProductsResp{}
	// 3. 执行ES查询
	// 在已有代码中找到ES查询构建部分，修改为：
	searchService := l.svcCtx.EsClient.Search().
		Index(biz.ProductEsIndexName).
		Query(boolQuery).
		From(from).
		Size(size)

	// 添加排序逻辑
	// 新增新品排序
	if in.New {
		searchService.SortBy(elastic.NewFieldSort("created_at").Desc()) // 按创建时间倒序
		searchService.SortBy(elastic.NewFieldSort("updated_at").Desc())
	}

	// 执行查询
	searchResult, err := searchService.Do(l.ctx)

	if err != nil || searchResult.Error != nil {
		l.Errorw("query es error", logx.Field("err", err))
		return res, nil
	}

	products := make([]*product.Product, searchResult.Hits.TotalHits.Value)
	for i, hit := range searchResult.Hits.Hits {
		var p product.Product
		bytes, err := hit.Source.MarshalJSON()
		if err != nil {
			l.Logger.Errorw("marshal json error")
			break
		}
		if err := json.Unmarshal(bytes, &p); err != nil {
			l.Logger.Errorw("unmarshal json error")
			break
		}
		products[i] = &p
	}
	res.Products = products
	return res, nil
}

// 辅助函数：转换字符串切片为interface切片
func stringSliceToInterface(s []string) []interface{} {
	r := make([]interface{}, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}
