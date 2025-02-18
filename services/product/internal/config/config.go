package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"time"
)

type Config struct {
	// gRPC 配置
	zrpc.RpcServerConf
	MysqlConfig         MysqlConfig
	RedisConf           redis.RedisConf
	ElasticsearchConfig ElasticsearchConfig
	QiNiu               QiNiu
	Consul              consul.Conf
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
