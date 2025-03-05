package audit

import (
	"context"
	"fmt"
	"jijizhazha1024/go-mall/services/audit/audit"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"jijizhazha1024/go-mall/common/consts/biz"
)

var auditClient audit.AuditClient
var auditOnce sync.Once

func setupAuditClient(t *testing.T) {
	auditOnce.Do(func() {
		conn, err := grpc.NewClient(
			fmt.Sprintf("127.0.0.1:%d", biz.AuditRpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			t.Fatalf("连接审计服务失败: %v", err)
		}
		auditClient = audit.NewAuditClient(conn)
	})
}

func TestCreateAuditLog(t *testing.T) {
	setupAuditClient(t)
	ctx := context.Background()
	now := time.Now().Unix()

	t.Run("正常创建审计日志", func(t *testing.T) {
		validReq := &audit.CreateAuditLogReq{
			UserId:      1001,
			ActionType:  "UPDATE",
			TargetTable: "products",
			TargetId:    2001,
			ClientIp:    "192.168.1.1",
			ServiceName: "product_service",
			// 以下为非必填字段
			ActionDescription: "更新商品价格",
			OldData:           `{"price": 100}`,
			NewData:           `{"price": 120}`,
			CreateAt:          now,
		}

		// 调用服务
		resp, err := auditClient.CreateAuditLog(ctx, validReq)
		assert.NoError(t, err)
		assert.Equal(t, uint32(0), resp.StatusCode)
		assert.True(t, resp.Ok)
	})
	t.Run("创建审计日志时缺少必填字段", func(t *testing.T) {
		invalidReq := &audit.CreateAuditLogReq{
			UserId:      1001,
			ActionType:  "UPDATE",
			TargetTable: "products",
			TargetId:    2001,
			ServiceName: "product_service",
			// 以下为非必填字段
			ActionDescription: "更新商品价格",
			OldData:           `{"price": 100}`,
			NewData:           `{"price": 120}`,
			CreateAt:          now,
		}
		// 调用服务
		_, err := auditClient.CreateAuditLog(ctx, invalidReq)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client_ip is required")
	})

}
