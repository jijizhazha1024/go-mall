package strategy

import (
	"context"
	"jijizhazha1024/go-mall/services/ai/internal/core/model"
	"jijizhazha1024/go-mall/services/ai/internal/core/vars"
)

type CommandStrategy interface {
	Parse(ctx context.Context, input string, userID int) (model.AST, error)
	Validate(ctx context.Context, ast model.AST) error
	Execute(ctx context.Context, ast model.AST) (interface{}, error)
	GetCommandType() vars.CommandType
}
