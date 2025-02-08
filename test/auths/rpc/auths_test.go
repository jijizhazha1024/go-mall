package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"jijizhazha1024/go-mall/common/consts/biz"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/auths/auths"
)

var client auths.AuthsClient
var once sync.Once

func setupGRPCConnection(t *testing.T) {
	once.Do(func() {
		conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%d", biz.AuthsRpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			t.Fatalf("Failed to connect to RPC server: %v", err)
		}
		client = auths.NewAuthsClient(conn)
	})
}

func TestAuthenticationLogic_Authentication(t *testing.T) {
	setupGRPCConnection(t)

	resp, err := client.GenerateToken(context.Background(), &auths.AuthGenReq{
		UserId:   1,
		Username: "test",
	})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	assert.Equal(t, uint32(0), resp.StatusCode)

	res, err := client.Authentication(context.Background(), &auths.AuthReq{
		Token: resp.AccessToken,
	})
	if err != nil {
		t.Fatalf("Authentication failed: %v", err)
	}
	assert.Equal(t, uint32(0), res.StatusCode)
	t.Log(res)
}

func TestAuthenticationLogic_GenerateToken(t *testing.T) {
	setupGRPCConnection(t)

	resp, err := client.GenerateToken(context.Background(), &auths.AuthGenReq{
		UserId:   1,
		Username: "test",
	})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	assert.Equal(t, uint32(0), resp.StatusCode)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)

	t.Log(resp)
}

func TestAuthenticationLogic_RenewToken(t *testing.T) {
	setupGRPCConnection(t)

	resp, err := client.GenerateToken(context.Background(), &auths.AuthGenReq{
		UserId:   1,
		Username: "test",
	})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	assert.Equal(t, uint32(0), resp.StatusCode)

	// 这里假设 token 的有效期为 10 秒，refresh token 的有效期为 30 分钟
	time.Sleep(time.Second * 11)

	res, err := client.Authentication(context.Background(), &auths.AuthReq{
		Token: resp.AccessToken,
	})
	if err != nil {
		t.Fatalf("Authentication failed: %v", err)
	}
	if res.StatusCode == code.AuthExpired {
		t.Logf("exprie token is %s", resp.AccessToken)

		renewResp, err := client.RenewToken(context.Background(), &auths.AuthRenewalReq{
			RefreshToken: resp.RefreshToken,
		})
		if err != nil {
			t.Fatalf("RenewToken failed: %v", err)
		}
		assert.Equal(t, uint32(0), renewResp.StatusCode)
		t.Logf("renew token is %s", renewResp.AccessToken)
	} else {
		assert.Equal(t, uint32(0), res.StatusCode)
		t.Log(res)
	}
}
