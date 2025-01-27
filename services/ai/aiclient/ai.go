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
	Request  = ai.Request
	Response = ai.Response

	Ai interface {
		Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
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

func (m *defaultAi) Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	client := ai.NewAiClient(m.cli.Conn())
	return client.Ping(ctx, in, opts...)
}
