package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/carts/carts"

	"jijizhazha1024/go-mall/apis/carts/internal/svc"
	"jijizhazha1024/go-mall/apis/carts/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCartItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCartItemLogic {
	return &DeleteCartItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCartItemLogic) DeleteCartItem(req *types.DeleteCartReq) (resp *types.DeleteCartResp, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.CartRpc.DeleteCartItem(l.ctx, &carts.CartItemRequest{
		Id: req.Id,
	})
	if err != nil {
		if res != nil && res.StatusCode != code.Success {
			// 处理用户级别info 错误
			return nil, errors.New(int(res.StatusCode), res.StatusMsg)
		}
		l.Logger.Errorf("call rpc DeleteCartItem failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}
	resp = &types.DeleteCartResp{
		Success: true,
	}
	return
}
