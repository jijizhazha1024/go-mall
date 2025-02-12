package updateaddress

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
	resp, err := users_client.UpdateAddress(context.Background(), &users.UpdateAddressRequest{

		RecipientName:   "djj",
		PhoneNumber:     "13800138000",
		Province:        "广东省",
		City:            "广州市",
		DetailedAddress: "天河区",
		IsDefault:       true,
		AddressId:       1,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("update success", resp)
	t.Log("updatesuccess", resp)
}
