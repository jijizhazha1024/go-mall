package checkout

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/checkout/checkout"
)

var checkoutClient checkout.CheckoutServiceClient

func init() {
	conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%d", biz.CheckoutRpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	checkoutClient = checkout.NewCheckoutServiceClient(conn)
}
