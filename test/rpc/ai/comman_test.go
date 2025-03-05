package aic

import (
	"context"
	"jijizhazha1024/go-mall/services/ai/aiclient"
	"testing"
)

func TestCommand(t *testing.T) {
	resp, err := aiClient.NLPExecutor(context.Background(), &aiclient.NLPExecutorReq{
		Command: "查询无线耳机，预算大概是100到500元,需要续航时间24小时，支持蓝牙",
		UserId:  1,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(resp.Data))
}
