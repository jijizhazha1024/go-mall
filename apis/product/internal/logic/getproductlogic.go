package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/services/product/product"

	"jijizhazha1024/go-mall/apis/product/internal/svc"
	"jijizhazha1024/go-mall/apis/product/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductLogic) GetProduct(req *types.GetProductReq) (resp *types.GetProductResp, err error) {
	res, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.GetProductReq{
		Id: uint32(req.ID),
	})
	if err != nil {
		l.Logger.Errorw("rpc get product detail  failed", logx.Field("err", err))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}
	resp = &types.GetProductResp{
		ID:          int64(res.Product.Id),
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Picture:     res.Product.Picture,
		Stock:       res.Product.Stock,
		Price:       float64(res.Product.Price),
		Categories:  res.Product.Categories,
	}
	return
}
