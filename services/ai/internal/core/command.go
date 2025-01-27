package core

import (
	"context"
	"jijizhazha1024/go-mall/services/ai/internal/config"
	"jijizhazha1024/go-mall/services/ai/internal/core/strategy"
	"jijizhazha1024/go-mall/services/ai/internal/core/vars"
)

type Command struct {
	strategies map[vars.CommandType]strategy.CommandStrategy
	conf       *config.Config
}

func NewCommand(ctx context.Context, conf *config.Config) *Command {
	return &Command{
		conf:       conf,
		strategies: make(map[vars.CommandType]strategy.CommandStrategy),
	}
}
func (c *Command) Handler(ctx context.Context, input string, userID int) (interface{}, error) {
	detectCommandType, err := c.detectCommandType(ctx, input)
	if err != nil {
		return nil, err
	}
	s := c.getStrategy(detectCommandType)
	ast, err := s.Parse(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	if err := s.Validate(ctx, userID, ast); err != nil {
		return nil, err
	}
	execute, err := s.Execute(ctx, userID, ast)
	if err != nil {
		return nil, err
	}
	return execute, nil
}

// parse user command type
func (c *Command) detectCommandType(ctx context.Context, input string) (vars.CommandType, error) {
	return vars.QueryProductCommand, nil
}
func (c *Command) register(s strategy.CommandStrategy) {
	c.strategies[s.GetCommandType()] = s
}
func (c *Command) getStrategy(tp vars.CommandType) strategy.CommandStrategy {
	if s, ok := c.strategies[tp]; ok {
		return s
	}
	return nil
}
