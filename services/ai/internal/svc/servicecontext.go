package svc

import (
	"jijizhazha1024/go-mall/services/ai/internal/config"
	"jijizhazha1024/go-mall/services/ai/internal/core"
)

type ServiceContext struct {
	Config  config.Config
	Command *core.Command
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Command: core.NewCommand(&c),
		Config:  c,
	}
}
