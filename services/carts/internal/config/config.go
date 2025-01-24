package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql MysqlConfig
}

type MysqlConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}
