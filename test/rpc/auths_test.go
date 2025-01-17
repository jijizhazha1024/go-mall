package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/auths/auths"
	"testing"
)

/*
由于go-zero进行编写测试模块不方便，这里使用原生的grpc进行测试
*/

var client auths.AuthsClient

func init() {
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", biz.AuthsRpcPort), grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = auths.NewAuthsClient(conn)
}
func TestAuthenticationLogic_Authentication(t *testing.T) {
	resp, err := client.Authentication(context.Background(), &auths.AuthReq{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
