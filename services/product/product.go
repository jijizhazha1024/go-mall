package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"jijizhazha1024/go-mall/common/consts/biz"
	"log"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"jijizhazha1024/go-mall/services/product/internal/config"
	"jijizhazha1024/go-mall/services/product/internal/logic"
	"jijizhazha1024/go-mall/services/product/internal/server"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"
)

var configFile = flag.String("f", "etc/product.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)
	// 创建一个 context 来控制定时任务的生命周期
	ctxWithCancel := context.Background()
	_, cancel := context.WithCancel(ctxWithCancel)
	// 启动定时任务扫描
	go func() {
		ticker := time.NewTicker(biz.ScanProductPVTime) // 每 5 小时执行一次
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				{
					err := logic.ScanHotProducts(ctx, ctxWithCancel)
					if err != nil {
						log.Println(err)
						return
					}
				}
			case <-ctxWithCancel.Done(): // 如果服务停止，退出定时任务
				log.Println("Stopping hot products scan task")
				return
			}
		}
	}()
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		product.RegisterProductCatalogServiceServer(grpcServer, server.NewProductCatalogServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// 注册服务
	if err := consul.RegisterService(c.ListenOn, c.Consul); err != nil {
		logx.Errorw("register service error", logx.Field("err", err))
		panic(err)
	}
	defer s.Stop()
	// 在服务停止时调用 cancel 函数，通知定时任务退出
	defer cancel()
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
