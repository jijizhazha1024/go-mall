package auths

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/auths/auths"

	"github.com/stretchr/testify/assert"
)

var client auths.AuthsClient
var once1 sync.Once
var clientIP string

func init() {
	// 获取客户端IP
	clientIP = "127.0.0.1"
}
func setupGRPCConnection(t *testing.T) {
	once1.Do(func() {
		conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%d", biz.AuthsRpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			t.Fatalf("Failed to connect to RPC server: %v", err)
		}
		client = auths.NewAuthsClient(conn)
	})
}

// 验证token
func TestAuthenticationLogic_Authentication(t *testing.T) {
	setupGRPCConnection(t)

	resp, err := client.GenerateToken(context.Background(), &auths.AuthGenReq{
		UserId:   4,
		Username: "test",
		ClientIp: clientIP,
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

// 签发token
func TestAuthenticationLogic_GenerateToken(t *testing.T) {
	setupGRPCConnection(t)

	resp, err := client.GenerateToken(context.Background(), &auths.AuthGenReq{
		UserId:   4,
		Username: "test",
		ClientIp: clientIP,
	})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	assert.Equal(t, uint32(0), resp.StatusCode)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)

	t.Log(resp)
}

// 续期token
func TestAuthenticationLogic_RenewToken(t *testing.T) {
	setupGRPCConnection(t)

	resp, err := client.GenerateToken(context.Background(), &auths.AuthGenReq{
		UserId:   4,
		Username: "test",
		ClientIp: clientIP,
	})
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	assert.Equal(t, uint32(0), resp.StatusCode)

	// 这里假设 token 的有效期为 10 秒，refresh token 的有效期为 30 分钟
	time.Sleep(time.Second * 11)

	res, err := client.Authentication(context.Background(), &auths.AuthReq{
		Token:    resp.AccessToken,
		ClientIp: clientIP,
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
