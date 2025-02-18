package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/product/product"

	"jijizhazha1024/go-mall/apis/product/internal/svc"
	"jijizhazha1024/go-mall/apis/product/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductListLogic {
	return &GetProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductListLogic) GetProductList(req *types.GetProductListReq) (resp *types.GetProductListResp, err error) {
	// 调用 RPC 服务获取分页商品列表
	res, err := l.svcCtx.ProductRpc.GetAllProduct(l.ctx, &product.GetAllProductsReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	// 处理 RPC 调用失败
	if err != nil {
		l.Logger.Errorw("failed to call GetProductList RPC",
			logx.Field("err", err),
			logx.Field("page", req.Page),
			logx.Field("page_size", req.PageSize))
		return nil, errors.New(code.ServerError, "RPC call failed")
	}

	// 处理 RPC 返回结果为空的情况
	if res == nil {
		l.Logger.Errorw("RPC GetProductList returned nil response",
			logx.Field("page", req.Page),
			logx.Field("page_size", req.PageSize))
		return nil, errors.New(code.ServerError, "RPC response is nil")
	}

	// 将 RPC 响应转换为 HTTP 响应
	products := make([]*types.Product, len(res.Products))
	for i, p := range res.Products {
		products[i] = &types.Product{
			ID:    int64(p.Id),
			Name:  p.Name,
			Stock: p.Stock,
			Price: float64(p.Price),
		}
	}

	// 构造 HTTP 响应
	resp = &types.GetProductListResp{
		Products: products,
		Total:    res.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	return resp, nil
}
