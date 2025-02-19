package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/apis/carts/internal/svc"
	"jijizhazha1024/go-mall/apis/carts/internal/types"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/carts/carts"
	"jijizhazha1024/go-mall/services/product/product"
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
	// 获取购物车信息
	userId := l.ctx.Value(biz.UserIDKey).(uint32)
	res, err := l.svcCtx.CartsRpc.CartItemList(l.ctx, &carts.UserInfo{
		Id: int32(userId),
	})

	// 处理 RPC 失败
	if err != nil {
		l.Logger.Errorw("call rpc CartItemList failed",
			logx.Field("err", err),
			logx.Field("user_id", userId))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}

	// 处理 RPC 返回结果为空的情况
	if res == nil {
		l.Logger.Errorw("rpc CartItemList returned nil response",
			logx.Field("user_id", userId))
		return nil, errors.New(code.ServerError, "RPC response is nil")
	}

	// 处理业务错误
	if res.StatusCode != code.Success {
		l.Logger.Debugw("rpc CartItemList returned business error",
			logx.Field("user_id", userId),
			logx.Field("status_code", res.StatusCode),
			logx.Field("status_msg", res.StatusMsg))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}

	// 获取商品ID列表
	ids := make([]int32, 0)
	for _, item := range res.Data {
		ids = append(ids, item.ProductId)
	}

	// 如果购物车为空，返回空商品列表
	if len(ids) == 0 {
		l.Logger.Infow("no product item in cart",
			logx.Field("user_id", userId),
			logx.Field("total", res.Total))
		return &types.CartItemListResp{
			Total:    res.Total,
			CartInfo: ConvertCartInfoResponse(res.Data),
		}, nil
	}

	// 获取商品详细信息
	productsList := make([]*types.CartInfoResponse, 0)
	for _, id := range ids {
		res_info, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.GetProductReq{
			Id: uint32(id),
		})
		if err != nil {
			l.Logger.Errorw("failed to get product", logx.Field("product_id", id), logx.Field("error", err))
			continue
		}

		// 获取商品详细信息
		if res_info != nil && res_info.Product != nil {
			product := res_info.Product

			// 处理每个购物车商品的详细信息
			for _, item := range res.Data {
				if item.ProductId == id {
					tmpCartInfo := &types.CartInfoResponse{
						Id:        item.Id,
						UserId:    item.UserId,
						ProductId: item.ProductId,
						Quantity:  item.Quantity,
						Product: []interface{}{
							map[string]interface{}{
								"product_id":    product.Id,
								"product_name":  product.Name,
								"product_image": product.Picture,
								"product_price": product.Price,
							},
						},
					}
					productsList = append(productsList, tmpCartInfo)
				}
			}
		}
	}

	// 正常返回
	l.Logger.Infow("Cart item list retrieved successfully",
		logx.Field("user_id", userId),
		logx.Field("total", res.Total))

	return &types.CartItemListResp{
		Total:    res.Total,
		CartInfo: productsList, // 返回包含详细信息的商品列表
	}, nil
}
