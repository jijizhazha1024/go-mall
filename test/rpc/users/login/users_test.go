package login

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/auths/auths"
	"jijizhazha1024/go-mall/services/users/users"
	"sync"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var users_client users.UsersClient
var auths_client auths.AuthsClient
var once1 sync.Once

func initusers() {
	once1.Do(func() {
		conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%d", biz.UsersRpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		conn1, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%d", biz.AuthsRpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)

		}
		users_client = users.NewUsersClient(conn)
		auths_client = auths.NewAuthsClient(conn1)
	})
}

func TestUsersRpc(t *testing.T) {
	initusers()
	resp, err := users_client.Login(context.Background(), &users.LoginRequest{
		Email:    "test9@test.com",
		Password: "1234567",
	})
	if err != nil {

		t.Fatal(err)
	}

	if resp.StatusCode == 0 {
		auths_res, err := auths_client.GenerateToken(context.Background(), &auths.AuthGenReq{
			UserId:   resp.UserId,
			Username: resp.UserName,
		})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("login success", resp, auths_res)
		t.Log("login success", resp)
	} else {
		fmt.Println("login failed", resp)
		t.Log("register failed", resp)
	}

}
