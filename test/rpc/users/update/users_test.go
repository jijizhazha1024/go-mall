package update

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/users/users"
	"sync"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var users_client users.UsersClient
var once1 sync.Once

func initusers() {
	once1.Do(func() {
		conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%d", biz.UsersRpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		users_client = users.NewUsersClient(conn)
	})
}

func TestUsersRpc(t *testing.T) {
	initusers()
	resp, err := users_client.UpdateUser(context.Background(), &users.UpdateUserRequest{

		UsrName: "test4",
		UserId:  4, //通过id修改
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("update success", resp)
	t.Log("update success", resp)
}
