package aic

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/ai/ai"
	"jijizhazha1024/go-mall/services/ai/aiclient"
)

var aiClient aiclient.Ai

func init() {
	conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%d", biz.AIRpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	aiClient = ai.NewAiClient(conn)
}
