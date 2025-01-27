// Code generated by goctl. DO NOT EDIT.
// Source: ai.proto

package server

import (
	"context"

	"jijizhazha1024/go-mall/services/ai/ai"
	"jijizhazha1024/go-mall/services/ai/internal/logic"
	"jijizhazha1024/go-mall/services/ai/internal/svc"
)

type AiServer struct {
	svcCtx *svc.ServiceContext
	ai.UnimplementedAiServer
}

func NewAiServer(svcCtx *svc.ServiceContext) *AiServer {
	return &AiServer{
		svcCtx: svcCtx,
	}
}

func (s *AiServer) Ping(ctx context.Context, in *ai.Request) (*ai.Response, error) {
	l := logic.NewPingLogic(ctx, s.svcCtx)
	return l.Ping(in)
}
