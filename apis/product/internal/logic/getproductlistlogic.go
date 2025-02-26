package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/common/consts/biz"
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
	userID, ok := l.ctx.Value(biz.UserIDKey).(uint32)
	if !ok {
		return nil, errors.New(code.AuthBlank, code.AuthBlankMsg)
	}
	var res *product.GetAllProductsResp
	if userID == 0 {
		// 调用 RPC 服务获取分页商品列表
		res, err = l.svcCtx.ProductRpc.GetAllProduct(l.ctx, &product.GetAllProductsReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		})
	} else {
		res, err = l.svcCtx.ProductRpc.RecommendProduct(l.ctx, &product.RecommendProductReq{
			UserId: int32(userID),
			Paginator: &product.RecommendProductReq_Paginator{
				Page:     req.Page,
				PageSize: req.PageSize,
			},
		})
	}
	// 处理 RPC 调用失败
	if err != nil {
		l.Logger.Errorw("call rpc ProductRpc failed",
			logx.Field("err", err),
			logx.Field("page", req.Page),
			logx.Field("page_size", req.PageSize))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	}
	if res.StatusCode != code.Success {
		// 可以记录日志
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}

	// 将 RPC 响应转换为 HTTP 响应
	products := make([]*types.Product, len(res.Products))
	for i, p := range res.Products {
		products[i] = &types.Product{
			ID:          int64(p.Id),
			Name:        p.Name,
			Stock:       p.Stock,
			Price:       p.Price,
			Picture:     p.Picture,
			Description: p.Description,
			Categories:  p.Categories,
			Sold:        p.Sold,
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
