package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	paymentM "jijizhazha1024/go-mall/dal/model/payment"
	"jijizhazha1024/go-mall/services/order/order"
	"log"
	"net/http"
	"time"

	"jijizhazha1024/go-mall/services/payment/internal/config"
	"jijizhazha1024/go-mall/services/payment/internal/server"
	"jijizhazha1024/go-mall/services/payment/internal/svc"
	"jijizhazha1024/go-mall/services/payment/payment"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/payment.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		payment.RegisterPaymentServer(grpcServer, server.NewPaymentServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	if err := consul.RegisterService(c.ListenOn, c.Consul); err != nil {
		logx.Errorw("register service error", logx.Field("err", err))
		panic(err)
	}
	http.HandleFunc("/notify", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		fmt.Println("notify  sdsdsds")
		// DecodeNotification 内部已调用 VerifySign 方法验证签名
		var noti, err = ctx.Alipay.DecodeNotification(request.Form)
		if err != nil {
			log.Println("Failed to decode notification:", err)
			return
		}
		order_Rpc := ctx.OrderRpc
		payment_model := ctx.PaymentModel

		// 根据通知状态处理业务逻辑
		switch noti.TradeStatus {
		case "TRADE_SUCCESS":
			// 解析时间字符串
			layout := "2006-01-02 15:04:05" // Go语言时间格式
			paymentTime, err := time.Parse(layout, noti.GmtPayment)
			timestamp := paymentTime.Unix()
			// 支付成功
			payment_model.UpdateInfoByOrderId(context.Background(), &paymentM.Payments{
				PreOrderId:    noti.OutTradeNo,
				OrderId:       sql.NullString{String: noti.OutTradeNo}, // 支付成功后更新
				TransactionId: sql.NullString{String: noti.TradeNo},
				Status:        int64(payment.PaymentStatus_PAYMENT_STATUS_PAID), // 初始状态：待支付
				UpdatedAt:     time.Now(),
				PaidAt:        sql.NullInt64{Int64: timestamp},
			})
			_, err = order_Rpc.UpdateOrder2Payment(context.Background(), &order.UpdateOrder2PaymentRequest{
				OrderId:     noti.OutTradeNo,
				UserId:      0,
				OrderStatus: order.OrderStatus_ORDER_STATUS_PAID,
			})
			if err != nil {
				return
			}
			log.Printf("Order %s payment succeeded, amount: %s", noti.OutTradeNo, noti.TotalAmount)
			// 在这里执行支付成功的业务逻辑，例如更新订单状态、记录日志等

		}
		// 返回确认响应给支付宝
		alipay.ACKNotification(writer)
	})
	// 启动 HTTP 服务
	go func() {
		if err := http.ListenAndServe(":12345", nil); err != nil {
			logx.Errorw("http server error", logx.Field("err", err))
		}
	}()
	// 启动RabbitMQ消费者
	go func() {
		consumeRabbitMQ(ctx)
	}()
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

// RabbitMQ消费者逻辑
func consumeRabbitMQ(ctx *svc.ServiceContext) {
	conn, err := amqp.Dial("amqp://admin:admin@124.71.72.124:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel: %v", err)
	}

	queueName := "delayed_queue"
	_, err = ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 排他性
		false,     // 等待确认
		nil,       // 队列参数
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		queueName, // 队列名称
		"",        // 消费者标签
		true,      // 自动确认（ack）
		false,     // 排他性
		false,     // 本地消息
		false,     // 等待确认
		nil,       // 参数
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	log.Printf("Consumer started. Waiting for messages from queue: %s", queueName)

	for msg := range msgs {
		orderID := string(msg.Body)
		log.Printf("Received order ID: %s", orderID)
		// 在这里处理订单ID，例如更新订单状态等
		paymentInfo, err := ctx.PaymentModel.FindOneByOrderId(context.Background(), orderID)
		if err != nil {
			return
		}
		if paymentInfo.Status == 1 {
			//调用oderRpc通知支付超时
		}

	}
}
