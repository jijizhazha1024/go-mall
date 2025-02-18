package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/services/product/product"

	"jijizhazha1024/go-mall/apis/product/internal/svc"
	"jijizhazha1024/go-mall/apis/product/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductByIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductByIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductByIDLogic {
	return &GetProductByIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductByIDLogic) GetProductByID(req *types.GetProductByIDReq) (resp *types.GetProductByIDResp, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.GetProductReq{
		Id: uint32(req.ID),
	})
	if err != nil {
		l.Logger.Errorf("rpc get product detail  failed", logx.Field("err", err))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}
	resp = &types.GetProductByIDResp{

		ID:          int64(res.Product.Id),
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Picture:     res.Product.Picture,
		Stock:       res.Product.Stock,
		Price:       float64(res.Product.Price),
		Categories:  res.Product.Categories,
	}

	return resp, err
}
