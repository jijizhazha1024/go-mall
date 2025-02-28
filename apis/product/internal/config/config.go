package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	ProductRpc     zrpc.RpcClientConf
	AuthsRpc       zrpc.RpcClientConf
	OptionPathList []string
}
