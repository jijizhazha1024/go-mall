package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/carts/carts"
	"jijizhazha1024/go-mall/services/product/product"

	"jijizhazha1024/go-mall/apis/carts/internal/svc"
	"jijizhazha1024/go-mall/apis/carts/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubCartItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubCartItemLogic {
	return &SubCartItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubCartItemLogic) SubCartItem(req *types.SubCartReq) (resp *types.SubCartResp, err error) {
	userId := l.ctx.Value(biz.UserIDKey).(uint32)

	// 1. 调用 GetProduct RPC 获取商品详情
	productRes, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.GetProductReq{
		Id: uint32(req.ProductId),
	})
	if err != nil {
		l.Logger.Errorw("call rpc GetProduct failed",
			logx.Field("err", err),
			logx.Field("product_id", req.ProductId))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}

	// 2. 检查商品是否存在
	if productRes == nil || productRes.Product == nil {
		l.Logger.Errorw("rpc GetProduct returned nil response",
			logx.Field("product_id", req.ProductId))
		return nil, errors.New(code.ProductInfoRetrievalFailed, code.ProductInfoRetrievalFailedMsg)
	}

	// 3. 调用 SubCartItem RPC 从购物车减少商品数量
	res, err := l.svcCtx.CartsRpc.SubCartItem(l.ctx, &carts.CartItemRequest{
		UserId:    int32(userId),
		ProductId: req.ProductId,
	})
	if err != nil {
		l.Logger.Errorw("call rpc SubCartItem failed",
			logx.Field("err", err),
			logx.Field("user_id", userId),
			logx.Field("product_id", req.ProductId))
		return nil, errors.New(code.CartSubFailed, code.CartSubFailedMsg)
	}

	// 4. 处理 RPC 返回 nil 的情况
	if res == nil {
		l.Logger.Errorw("rpc SubCartItem returned nil response",
			logx.Field("user_id", userId),
			logx.Field("product_id", req.ProductId))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}

	// 5. 处理业务错误
	if res.StatusCode != code.Success {
		l.Logger.Debugw("rpc SubCartItem returned business error",
			logx.Field("user_id", userId),
			logx.Field("product_id", req.ProductId),
			logx.Field("status_code", res.StatusCode),
			logx.Field("status_msg", res.StatusMsg))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}

	// 6. 记录成功日志并返回结果
	l.Logger.Infow("Cart item subtracted successfully",
		logx.Field("user_id", userId),
		logx.Field("product_id", req.ProductId))

	return &types.SubCartResp{
		Id: res.Id,
	}, nil
}
