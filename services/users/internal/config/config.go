package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	MysqlConfig MysqlConfig
	HttpClient  HttpClientConf
	AuditRpc    zrpc.RpcClientConf
	Consul      consul.Conf
	Cache       cache.CacheConf
	RedisConf   redis.RedisConf
}
type MysqlConfig struct {
	DataSource  string
	Conntimeout int
}
type HttpClientConf struct {
	Timeout      int
	MaxIdleConns int
}
