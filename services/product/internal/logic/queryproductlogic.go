package logic

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"

	"github.com/zeromicro/go-zero/core/logx"
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

// 根据条件查询商品
func (l *QueryProductLogic) QueryProduct(in *product.QueryProductReq) (*product.GetAllProductsResp, error) {
	// 1. 创建ES查询客户端
	client := l.svcCtx.EsClient
	indexName := biz.ProductEsIndexName // 替换为实际的ES索引名

	// 2. 构建基础查询
	boolQuery := elastic.NewBoolQuery()
	// 3. 条件组装
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

	// 4. 分页处理
	from := 0
	size := 10 // 默认分页大小
	if in.Paginator != nil {
		if in.Paginator.PageSize > 0 {
			size = int(in.Paginator.PageSize)
		}
		if in.Paginator.Page > 1 {
			from = (int(in.Paginator.Page) - 1) * size
		}
	}

	// 5. 执行ES查询
	searchResult, err := client.Search().
		Index(indexName).
		Query(boolQuery).
		From(from).
		Size(size).
		Do(l.ctx)

	if err != nil || searchResult.Error != nil {
		l.Errorw("query es error", logx.Field("err", err))
		return nil, err
	}

	var products []*product.Product
	for _, hit := range searchResult.Hits.Hits {
		var p product.Product
		bytes, err := hit.Source.MarshalJSON()
		if err != nil {
			l.Logger.Errorw("marshal json error")
			return nil, err
		}
		if err := json.Unmarshal(bytes, &p); err != nil {
			l.Logger.Errorw("unmarshal json error")
			return nil, err
		}
		products = append(products, &p)
	}
	return &product.GetAllProductsResp{
		Products: products,
	}, nil
}

// 辅助函数：转换字符串切片为interface切片
func stringSliceToInterface(s []string) []interface{} {
	r := make([]interface{}, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}
