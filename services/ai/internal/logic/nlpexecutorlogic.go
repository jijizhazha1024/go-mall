package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/ai/ai"
	"jijizhazha1024/go-mall/services/ai/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type NLPExecutorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNLPExecutorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NLPExecutorLogic {
	return &NLPExecutorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *NLPExecutorLogic) NLPExecutor(in *ai.NLPExecutorReq) (*ai.NLPExecutorResp, error) {
	// todo: add your logic here and delete this line

	return &ai.NLPExecutorResp{}, nil
}
