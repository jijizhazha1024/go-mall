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
	boolQuery := l.buildESQuery(in)
	// 分页参数
	pageSize := biz.DefaultPageSize // 默认分页大小
	from := 0
	if in.Paginator != nil && in.Paginator.PageSize > 0 {
		pageSize = int(in.Paginator.PageSize)
	}
	if in.Paginator != nil && in.Paginator.Page > 0 {
		from = (int(in.Paginator.Page) - 1) * pageSize
	}
	// 构建搜索服务
	searchService := l.svcCtx.EsClient.Search().
		Index(biz.ProductEsIndexName).
		Query(boolQuery).
		From(from).
		Size(pageSize)

	// 新增排序逻辑
	if in.New {
		searchService.SortBy(
			elastic.NewFieldSort("created_at").Desc(),
			elastic.NewFieldSort("updated_at").Desc(),
		)
	}
	if in.Hot {
		searchService.SortBy(
			elastic.NewFieldSort("price").Asc(),
			elastic.NewScoreSort().Desc(),
		)
	}
	searchResult, err := searchService.Do(l.ctx)
	if err != nil {
		logx.Errorw("elasticsearch query error", logx.Field("err", err))
		return nil, err
	}
	// 处理查询结果
	var products []*product.Product
	for _, hit := range searchResult.Hits.Hits {
		var p *product.Product
		if err := json.Unmarshal(hit.Source, &p); err != nil {
			continue
		}
		products = append(products, p)
	}

	return &product.GetAllProductsResp{
		Total:    searchResult.TotalHits(),
		Products: products,
	}, nil
}
func (l *QueryProductLogic) buildESQuery(req *product.QueryProductReq) *elastic.BoolQuery {
	boolQuery := elastic.NewBoolQuery()

	// 商品名称模糊匹配
	if req.Name != "" {
		boolQuery.Filter(elastic.NewMatchQuery("name", req.Name))
	}
	if req.Keyword != "" {
		// 修改后代码
		keywordBool := elastic.NewBoolQuery()
		keywordBool.Should(
			elastic.NewMultiMatchQuery(req.Keyword, "name^1", "description^2"),
			elastic.NewMatchPhraseQuery("name", req.Keyword).Boost(1),
			elastic.NewMatchPhraseQuery("description", req.Keyword).Boost(3),
			elastic.NewWildcardQuery("description.keyword", "*"+req.Keyword+"*").Boost(2), // 新增通配符查询
			elastic.NewTermQuery("description.keyword", req.Keyword).Boost(5),
		)
		keywordBool.MinimumNumberShouldMatch(1)
		boolQuery.Must(keywordBool) // 从Filter改为Must以保留评分
	}

	// 分类筛选（数组匹配）
	if len(req.Category) > 0 {
		termsQuery := elastic.NewTermsQuery("category.keyword", stringSliceToInterface(req.Category)...)
		boolQuery.Should(termsQuery)
	}

	// 价格区间筛选
	if req.Price != nil {
		rangeQuery := elastic.NewRangeQuery("price")
		if req.Price.Min > 0 {
			rangeQuery.Gte(req.Price.Min)
		}
		if req.Price.Max > 0 {
			rangeQuery.Lte(req.Price.Max)
		}
		boolQuery.Filter(rangeQuery)
	}

	return boolQuery
}

// 确保存在该转换函数（建议放在公共工具类中）
func stringSliceToInterface(s []string) []interface{} {
	r := make([]interface{}, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}
