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

	// 3. 检查库存是否足够
	if productRes.Product.Stock == 0 {
		l.Logger.Errorw("insufficient stock",
			logx.Field("product_id", req.ProductId),
			logx.Field("available_stock", productRes.Product.Stock))
		return nil, errors.New(code.InsufficientInventoryOfProduct, code.InsufficientInventoryOfProductMsg)
	}

	// 4. 调用 CreateCartItem RPC 添加到购物车
	res, err := l.svcCtx.CartsRpc.CreateCartItem(l.ctx, &carts.CartItemRequest{
		UserId:    int32(userId),
		ProductId: req.ProductId,
	})
	if err != nil {
		l.Logger.Errorw("call rpc CreateCartItem failed",
			logx.Field("err", err),
			logx.Field("user_id", userId),
			logx.Field("product_id", req.ProductId))
		return nil, errors.New(code.CartCreationFailed, code.CartCreationFailedMsg)
	}

	// 5. 处理 RPC 返回 nil 的情况
	if res == nil {
		l.Logger.Errorw("rpc CreateCartItem returned nil response",
			logx.Field("user_id", userId),
			logx.Field("product_id", req.ProductId))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}

	// 6. 处理业务错误
	if res.StatusCode != code.Success {
		l.Logger.Debugw("rpc CreateCartItem returned business error",
			logx.Field("user_id", userId),
			logx.Field("product_id", req.ProductId),
			logx.Field("status_code", res.StatusCode),
			logx.Field("status_msg", res.StatusMsg))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}

	// 7. 记录成功日志并返回结果
	l.Logger.Infow("Cart item created successfully",
		logx.Field("user_id", userId),
		logx.Field("product_id", req.ProductId))

	return &types.CreateCartResp{
		Id: res.Id,
	}, nil
}
