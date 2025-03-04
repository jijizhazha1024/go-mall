package main

import (
	"flag"
	"fmt"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"

	"jijizhazha1024/go-mall/apis/payment/internal/config"
	"jijizhazha1024/go-mall/apis/payment/internal/handler"
	"jijizhazha1024/go-mall/apis/payment/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/payment-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
