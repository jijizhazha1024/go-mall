package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"jijizhazha1024/go-mall/common/config"
	"time"
)

type Config struct {
	// gRPC 配置
	zrpc.RpcServerConf
	MysqlConfig   MysqlConfig
	RedisConf     redis.RedisConf
	ElasticSearch config.ElasticSearchConfig
	QiNiu         QiNiu
	Consul        consul.Conf
	InventoryRpc  zrpc.RpcClientConf
}
type MysqlConfig struct {
	DataSource  string
	Conntimeout int
}
type ElasticsearchConfig struct {
	Addresses             []string      `yaml:"addresses"`
	MaxIdleConnsPerHost   int           `yaml:"max_idle_conns_per_host"`
	ResponseHeaderTimeout time.Duration `yaml:"response_header_timeout"`
}
type QiNiu struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Domain    string
}
