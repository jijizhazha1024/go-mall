package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/apis/carts/internal/svc"
	"jijizhazha1024/go-mall/apis/carts/internal/types"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/carts/carts"
)

type CartItemListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartItemListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartItemListLogic {
	return &CartItemListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartItemListLogic) CartItemList(req *types.UserInfo) (resp *types.CartItemListResp, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.CartRpc.CartItemList(l.ctx, &carts.UserInfo{
		Id: req.Id,
	})

	// 处理 RPC 失败
	if err != nil {
		l.Logger.Errorw("call rpc CartItemList failed",
			logx.Field("err", err),
			logx.Field("user_id", req.Id))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}

	// 处理 RPC 返回结果为空的情况
	if res == nil {
		l.Logger.Errorw("rpc CartItemList returned nil response",
			logx.Field("user_id", req.Id))
		return nil, errors.New(code.ServerError, "RPC response is nil")
	}

	// 处理业务错误
	if res.StatusCode != code.Success {
		l.Logger.Debugw("rpc CartItemList returned business error",
			logx.Field("user_id", req.Id),
			logx.Field("status_code", res.StatusCode),
			logx.Field("status_msg", res.StatusMsg))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}

	// 正常返回
	l.Logger.Infow("Cart item list retrieved successfully",
		logx.Field("user_id", req.Id),
		logx.Field("total", res.Total))

	return &types.CartItemListResp{
		Total:    res.Total,
		CartInfo: ConvertCartInfoResponse(res.Data),
	}, nil
}
