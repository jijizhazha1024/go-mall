// Code generated by goctl. DO NOT EDIT.
// Source: ai.proto

package aiclient

import (
	"context"

	"jijizhazha1024/go-mall/services/ai/ai"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	NLPExecutorReq  = ai.NLPExecutorReq
	NLPExecutorResp = ai.NLPExecutorResp

	Ai interface {
		NLPExecutor(ctx context.Context, in *NLPExecutorReq, opts ...grpc.CallOption) (*NLPExecutorResp, error)
	}

	defaultAi struct {
		cli zrpc.Client
	}
)

func NewAi(cli zrpc.Client) Ai {
	return &defaultAi{
		cli: cli,
	}
}

func (m *defaultAi) NLPExecutor(ctx context.Context, in *NLPExecutorReq, opts ...grpc.CallOption) (*NLPExecutorResp, error) {
	client := ai.NewAiClient(m.cli.Conn())
	return client.NLPExecutor(ctx, in, opts...)
}
