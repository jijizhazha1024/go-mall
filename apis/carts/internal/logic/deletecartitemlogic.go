package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/biz"
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
	userId := l.ctx.Value(biz.UserIDKey).(uint32)
	res, err := l.svcCtx.CartsRpc.DeleteCartItem(l.ctx, &carts.CartItemRequest{
		UserId:    int32(userId),
		ProductId: req.ProductId,
	})

	// 处理 RPC 层面的错误
	if err != nil {
		l.Logger.Errorw("call rpc DeleteCartItem failed",
			logx.Field("err", err),
			logx.Field("user_id", req.UserId),
			logx.Field("product_id", req.ProductId))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}

	// 处理业务错误（RPC 可能返回非成功状态）
	if res == nil {
		l.Logger.Errorw("rpc DeleteCartItem returned nil response",
			logx.Field("request_id", req.Id))
		return nil, errors.New(code.ServerError, "RPC response is nil")
	}
	if res.StatusCode != code.Success {
		l.Logger.Debugw("rpc DeleteCartItem returned business error",
			logx.Field("request_id", req.Id),
			logx.Field("status_code", res.StatusCode),
			logx.Field("status_msg", res.StatusMsg))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}

	// 操作成功
	l.Logger.Infow("Cart item deleted successfully",
		logx.Field("request_id", req.Id),
		logx.Field("user_id", req.UserId),
		logx.Field("product_id", req.ProductId))

	return &types.DeleteCartResp{Success: true}, nil
}
