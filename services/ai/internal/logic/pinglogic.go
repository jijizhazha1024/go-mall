package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/ai/ai"
	"jijizhazha1024/go-mall/services/ai/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *ai.Request) (*ai.Response, error) {
	// todo: add your logic here and delete this line

	return &ai.Response{}, nil
}
