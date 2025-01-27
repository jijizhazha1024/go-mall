package strategy

import (
	"context"
	"jijizhazha1024/go-mall/services/ai/internal/core/model"
	"jijizhazha1024/go-mall/services/ai/internal/core/vars"
)

type ProductQueryStrategy struct{}

func (s *ProductQueryStrategy) Parse(ctx context.Context, input string, userID int) (model.AST, error) {
	// 示例：记录解析日志

	return &model.ProductQueryAST{
		Command:  "product_query",
		Page:     1,
		PageSize: 10,
		UserID:   userID,
	}, nil
}

func (s *ProductQueryStrategy) Validate(ctx context.Context, userID int, ast model.AST) error {
	return nil
}

func (s *ProductQueryStrategy) Execute(ctx context.Context, userID int, ast model.AST) (interface{}, error) {
	return model.ProductQueryResponseData{
		Products: []model.Product{
			{ID: "p1", Name: "商品A", Price: 100, Stock: 50},
		},
		Total: 1,
	}, nil
}
func (s *ProductQueryStrategy) GetCommandType() vars.CommandType {
	return vars.QueryProductCommand
}
