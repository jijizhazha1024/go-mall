package aic

import (
	"context"
	"jijizhazha1024/go-mall/services/ai/aiclient"
	"testing"
)

func TestCommand(t *testing.T) {
	resp, err := aiClient.NLPExecutor(context.Background(), &aiclient.NLPExecutorReq{
		Command: "查询智能手机，预算大概是4000到6000元,需要6.5英寸AMOLED屏幕，128GB存储",
		UserId:  1,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Data)
}
