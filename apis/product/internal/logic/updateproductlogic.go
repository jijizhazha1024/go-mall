package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/services/product/product"

	"jijizhazha1024/go-mall/apis/product/internal/svc"
	"jijizhazha1024/go-mall/apis/product/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductReq) (resp *types.UpdateProductResp, err error) {
	res, err := l.svcCtx.ProductRpc.UpdateProduct(l.ctx, &product.UpdateProductReq{
		Product: &product.Product{
			Id:          uint32(req.ID),
			Name:        req.Name,
			Description: req.Description,
			Picture:     req.Picture,
			Price:       float32(req.Price),
			Stock:       req.Stock,
			Categories:  req.Categories,
		},
	})
	if err != nil {
		l.Logger.Errorw("rpc update product  failed", logx.Field("err", err))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}
	resp = &types.UpdateProductResp{
		ProductID: req.ID,
	}
	return
}
