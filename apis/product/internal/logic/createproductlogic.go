package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/services/product/product"

	"jijizhazha1024/go-mall/apis/product/internal/svc"
	"jijizhazha1024/go-mall/apis/product/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (resp *types.CreateProductResp, err error) {
	res, err := l.svcCtx.ProductRpc.CreateProduct(l.ctx, &product.CreateProductReq{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       float32(req.Price),
		Stock:       req.Stock,
		Categories:  req.Categories,
	})
	if err != nil {
		l.Logger.Errorw("rpc create product  failed", logx.Field("err", err))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}
	resp = &types.CreateProductResp{
		ProductID: res.I,
	}
	return
}
