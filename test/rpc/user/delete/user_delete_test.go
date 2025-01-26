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
	conn, err := grpc.Dial(fmt.Sprintf("0.0.0.0:%d", biz.UsersRpcPort), grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	users_client = users.NewUsersClient(conn)
}

func TestUsersDeleteRpc(t *testing.T) {
	initusers()
	resp, err := users_client.DeleteUser(context.Background(), &users.DeleteUserRequest{
		UserId: 6,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("delete user success", resp)
	t.Log("delete user success", resp)

}
