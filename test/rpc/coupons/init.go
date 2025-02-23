package coupons

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/coupons/coupons"
)

var couponsClient coupons.CouponsClient

func init() {
	conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%d", biz.CouponsRpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	couponsClient = coupons.NewCouponsClient(conn)
}
