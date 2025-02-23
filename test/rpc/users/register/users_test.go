package register

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/users/users"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var users_client users.UsersClient

func initusers() {

	conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%d", biz.UsersRpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	users_client = users.NewUsersClient(conn)
}

func TestUsersRpc(t *testing.T) {
	initusers()
	resp, err := users_client.Register(context.Background(), &users.RegisterRequest{
		Email:           "djj2@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("register success", resp)
	t.Log("register success", resp)
}
