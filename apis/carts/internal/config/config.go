package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	AuthsRpc       zrpc.RpcClientConf
	CartsRpc       zrpc.RpcClientConf
	ProductRpc     zrpc.RpcClientConf
	WhitePathList  []string `json:",optional"`
	OptionPathList []string `json:",optional"`
}
