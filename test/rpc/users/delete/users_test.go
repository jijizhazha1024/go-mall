package delete

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

	//这里可以从token中获取user——id
	resp, err := users_client.DeleteUser(context.Background(), &users.DeleteUserRequest{
		UserId: 4,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("DELETE success", resp)
	t.Log("DELETE success", resp)
}
