package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/services/ai/ai"
	"jijizhazha1024/go-mall/services/ai/internal/core"
	"jijizhazha1024/go-mall/services/ai/internal/svc"
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

	res := &ai.NLPExecutorResp{}
	if in.UserId == 0 || in.Command == "" {
		return nil, nil
	}
	// 调用解析器
	command := core.NewCommand(&l.svcCtx.Config)
	handler, err := command.Handler(l.ctx, in.Command, int(in.UserId)) // any,error
	if err != nil {
		return res, nil
	}
	marshal, err := json.Marshal(handler)
	if err != nil {
		return nil, err
	}
	res.Data = marshal
	return res, nil
}
