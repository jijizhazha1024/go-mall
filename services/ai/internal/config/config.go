package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Gpt struct {
		ApiKey  string
		ModelID string
	}
}
