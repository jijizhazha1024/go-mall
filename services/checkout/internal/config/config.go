package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	MysqlConfig MysqlConfig
	RedisConf   redis.RedisConf
	Consul      consul.Conf

	InventoryRpc zrpc.RpcClientConf
	CouponsRpc   zrpc.RpcClientConf
	ProductRpc   zrpc.RpcClientConf
}

type MysqlConfig struct {
	DataSource  string
	Conntimeout int
}
