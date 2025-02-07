package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/common/config"
)

type Config struct {
	zrpc.RpcServerConf
	Gpt config.GptConfig
}
