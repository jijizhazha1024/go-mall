package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Gpt struct {
		ApiKey  string
		ModelID string
	}
	Consul     consul.Conf
	ProductRpc zrpc.RpcClientConf
}
