package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	MysqlConfig MysqlConfig
	Consul      consul.Conf
}

type MysqlConfig struct {
	DataSource  string
	Conntimeout int
}
