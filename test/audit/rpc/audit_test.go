package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/audit/audit"
	"sync"
	"testing"
	"time"
)

var auditRpc audit.AuditClient
var auditRpcOnce sync.Once

func setupAuditRpcServer(t *testing.T) {
	auditRpcOnce.Do(func() {
		conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%d", biz.AuditRpcPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			t.Fatalf("Failed to connect to RPC server: %v", err)
		}
		auditRpc = audit.NewAuditClient(conn)
	})
}

func TestCreateAuditLog(t *testing.T) {
	setupAuditRpcServer(t)
	data, err := json.Marshal(map[string]string{
		"user_id":  "1",
		"username": "test",
	})
	if err != nil {
		t.Fatalf("Failed to call CreateAuditLog: %v", err)
	}
	res, err := auditRpc.CreateAuditLog(context.Background(), &audit.CreateAuditLogReq{
		UserId:            1,
		Username:          "test",
		ActionType:        "test",
		ActionDescription: "test",
		TargetTable:       "test",
		TargetId:          1,
		OldData:           string(data),
		NewData:           string(data),
		CreateAt:          time.Now().Unix(),
	})
	if err != nil {
		t.Fatalf("Failed to call CreateAuditLog: %v", err)
	}

	fmt.Println(res)
}
