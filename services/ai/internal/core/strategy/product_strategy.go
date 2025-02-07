package strategy

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/common/utils/gpt"
	"jijizhazha1024/go-mall/services/ai/internal/config"
	"jijizhazha1024/go-mall/services/ai/internal/core/model"
	"jijizhazha1024/go-mall/services/ai/internal/core/prompt/product_query"
	"jijizhazha1024/go-mall/services/ai/internal/core/vars"
)

type ProductQueryStrategy struct {
	gpt *gpt.Gpt
}

func (s *ProductQueryStrategy) Parse(ctx context.Context, command string, userID int) (model.AST, error) {
	str, err := s.gpt.ChatWithModel(ctx, product_query.Prompt, fmt.Sprintf("用户输入:  ```%s``` ", command))
	if err != nil {
		logx.Infow("failed to parse ast", logx.Field("error", err), logx.Field("command", command))
		// 一般是余额，超时问题
		return nil, err
	}
	var ast = new(model.ProductQueryAST)
	if err := jsonx.UnmarshalFromString(str, ast); err != nil {
		// 解析错误
		logx.Infow("failed to unmarshal ast", logx.Field("error", err),
			logx.Field("command", command), logx.Field("ast", str),
		)
		return nil, vars.ErrASTParseFailed
	}
	ast.UserID = userID
	ast.Command = command
	return ast, nil
}

func (s *ProductQueryStrategy) Validate(ctx context.Context, ast model.AST) error {
	return nil
}

func (s *ProductQueryStrategy) Execute(ctx context.Context, ast model.AST) (interface{}, error) {
	// 调用rpc服务
	return nil, nil
}
func (s *ProductQueryStrategy) GetCommandType() vars.CommandType {
	return vars.QueryProductCommand
}
func NewProductQueryStrategy(conf *config.Config) *ProductQueryStrategy {
	strategy := &ProductQueryStrategy{
		gpt: gpt.NewGpt(conf.Gpt.ApiKey, conf.Gpt.ModelID),
	}
	return strategy
}
