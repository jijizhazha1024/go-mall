package core

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/services/ai/internal/config"
	"jijizhazha1024/go-mall/services/ai/internal/core/strategy"
	"jijizhazha1024/go-mall/services/ai/internal/core/vars"
)

type Command struct {
	factory strategy.StrategyFactory
	conf    *config.Config
}

func NewCommand(conf *config.Config) *Command {
	return &Command{
		factory: strategy.NewDefaultStrategyFactory(conf),
		conf:    conf,
	}
}
func (c *Command) Handler(ctx context.Context, input string, userID int) (interface{}, error) {
	commandType, err := c.detectCommandType(ctx, input)
	if err != nil {
		return nil, err
	}
	s := c.factory.CreateStrategy(commandType)
	if s == nil {
		return nil, vars.ErrUnknownCommand
	}
	ast, err := s.Parse(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	fmt.Println(ast)
	if err := s.Validate(ctx, ast); err != nil {
		return nil, err
	}
	execute, err := s.Execute(ctx, ast)
	if err != nil {
		return nil, err
	}
	return execute, nil
}

// parse user command type
func (c *Command) detectCommandType(ctx context.Context, input string) (vars.CommandType, error) {
	return vars.QueryProductCommand, nil
}
