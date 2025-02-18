package addaddress

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
	resp, err := users_client.AddAddress(context.Background(), &users.AddAddressRequest{

		RecipientName:   "test2",
		PhoneNumber:     "13803423123",
		Province:        "山东省",
		City:            "威海市",
		DetailedAddress: "环翠区",
		IsDefault:       true,
		UserId:          1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("add success", resp)
	t.Log("addsuccess", resp)
}
