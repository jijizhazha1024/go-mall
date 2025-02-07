package strategy

import (
	"jijizhazha1024/go-mall/services/ai/internal/config"
	"jijizhazha1024/go-mall/services/ai/internal/core/vars"
)

type StrategyFactory interface {
	CreateStrategy(commandType vars.CommandType) CommandStrategy
}

type DefaultStrategyFactory struct {
	conf *config.Config
}

func NewDefaultStrategyFactory(conf *config.Config) StrategyFactory {
	return &DefaultStrategyFactory{conf: conf}
}

func (f *DefaultStrategyFactory) CreateStrategy(commandType vars.CommandType) CommandStrategy {
	switch commandType {
	case vars.QueryProductCommand:
		return NewProductQueryStrategy(f.conf)
	default:
		return nil
	}
}
