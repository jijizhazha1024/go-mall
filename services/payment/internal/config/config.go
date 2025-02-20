package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"jijizhazha1024/go-mall/common/config"
)

type Config struct {
	zrpc.RpcServerConf
	Consul      consul.Conf
	MysqlConfig config.MysqlConfig
	RedisConf   redis.RedisConf
	Alipay      struct {
		AppId           string
		PrivateKey      string
		AlipayPublicKey string
		NotifyURL       string
		ReturnURL       string
	} `yaml:"alipay"`
}
