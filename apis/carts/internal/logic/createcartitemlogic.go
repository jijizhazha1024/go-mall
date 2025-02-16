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

type CreateCartItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCartItemLogic {
	return &CreateCartItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCartItemLogic) CreateCartItem(req *types.CreateCartReq) (resp *types.CreateCartResp, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.CartRpc.CreateCartItem(l.ctx, &carts.CartItemRequest{
		UserId:       req.UserId,
		ProductId:    req.ProductId,
		ProductName:  req.ProductName,
		ProductImage: req.ProuctImage,
		ProductPrice: req.ProductPrice,
		Quantity:     req.Quantity,
		Checked:      req.Checked,
	})

	// 处理 RPC 失败
	if err != nil {
		l.Logger.Errorw("call rpc CreateCartItem failed",
			logx.Field("err", err),
			logx.Field("user_id", req.UserId),
			logx.Field("product_id", req.ProductId))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}

	// 处理 RPC 返回 nil 的情况
	if res == nil {
		l.Logger.Errorw("rpc CreateCartItem returned nil response",
			logx.Field("user_id", req.UserId),
			logx.Field("product_id", req.ProductId))
		return nil, errors.New(code.ServerError, "RPC response is nil")
	}

	// 处理业务错误
	if res.StatusCode != code.Success {
		l.Logger.Debugw("rpc CreateCartItem returned business error",
			logx.Field("user_id", req.UserId),
			logx.Field("product_id", req.ProductId),
			logx.Field("status_code", res.StatusCode),
			logx.Field("status_msg", res.StatusMsg))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}

	// 操作成功
	l.Logger.Infow("Cart item created successfully",
		logx.Field("user_id", req.UserId),
		logx.Field("product_id", req.ProductId),
		logx.Field("cart_id", res.Id))

	return &types.CreateCartResp{
		Id: res.Id,
	}, nil
}
