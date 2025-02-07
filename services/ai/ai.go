package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"jijizhazha1024/go-mall/services/ai/ai"
	"jijizhazha1024/go-mall/services/ai/internal/config"
	"jijizhazha1024/go-mall/services/ai/internal/server"
	"jijizhazha1024/go-mall/services/ai/internal/svc"
	"os"
)

var configFile = flag.String("f", "etc/ai.yaml", "the config file")

func main() {

	flag.Parse()
	os.Setenv("GPT_API_KEY", "5b5ab09c-7298-40d7-b60e-433d21314f36")
	os.Setenv("GPT_MODEL_ID", "ep-20241002090911-md25k")
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		ai.RegisterAiServer(grpcServer, server.NewAiServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
