package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/ai/ai"

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
	userId, ok := l.ctx.Value(biz.UserIDKey).(uint32)
	resp = &types.Response{}
	if !ok {
		return nil, errors.New(code.AuthBlank, code.AuthBlankMsg)
	}

	executor, err := l.svcCtx.AiRpc.NLPExecutor(l.ctx, &ai.NLPExecutorReq{
		UserId:  userId,
		Command: req.Command,
	})
	if err != nil {
		logx.Errorw("call rpc NLPExecutor failed", logx.Field("err", err))
		return nil, err
	}
	if executor.StatusCode != code.Success {
		return nil, errors.New(int(executor.StatusCode), executor.StatusMsg)
	}
	if err := jsonx.Unmarshal(executor.Data, &resp.Result); err != nil {
		return resp, nil
	}
	return
}
