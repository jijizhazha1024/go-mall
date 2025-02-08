package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"jijizhazha1024/go-mall/common/config"
)

type Config struct {
	zrpc.RpcServerConf
	Consul        consul.Conf
	RabbitMQ      config.RabbitMQConfig
	Mysql         config.MysqlConfig
	ElasticSearch config.ElasticSearchConfig
}
