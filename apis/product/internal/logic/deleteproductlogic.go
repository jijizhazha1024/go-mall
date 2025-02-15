package logic

import (
	"context"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/services/product/product"

	"jijizhazha1024/go-mall/apis/product/internal/svc"
	"jijizhazha1024/go-mall/apis/product/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProductLogic) DeleteProduct(req *types.DeleteProductReq) (resp *types.DeleteProductResp, err error) {
	res, err := l.svcCtx.ProductRpc.DeleteProduct(l.ctx, &product.DeleteProductReq{
		Id: req.ID,
	})
	if err != nil {
		return &types.DeleteProductResp{
			Success: false,
		}, errors.New(int(res.StatusCode), res.StatusMsg)
	}
	resp = &types.DeleteProductResp{
		Success: true,
	}

	return
}
