package logic

import (
	"context"

	"jijizhazha1024/go-mall/apis/ai/internal/svc"
	"jijizhazha1024/go-mall/apis/ai/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiLogic {
	return &AiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiLogic) Ai(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
