package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/dal/model/audit"
	"jijizhazha1024/go-mall/services/audit/internal/config"
	"jijizhazha1024/go-mall/services/audit/internal/mq"
)

type ServiceContext struct {
	Config     config.Config
	AuditMQ    *mq.AuditMQ
	AuditModel audit.AuditModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	auditMq, err := mq.Init(c)
	if err != nil {
		logx.Error(err)
		panic(err)
	}
	return &ServiceContext{
		Config:  c,
		AuditMQ: auditMq,
	}
}
