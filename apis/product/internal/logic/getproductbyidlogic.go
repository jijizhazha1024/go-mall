package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
	"jijizhazha1024/go-mall/apis/product/internal/svc"
	"jijizhazha1024/go-mall/apis/product/internal/types"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/product/product"
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
	userID, ok := l.ctx.Value(biz.UserIDKey).(uint32)
	if !ok {
		return nil, errors.New(code.AuthBlank, code.AuthBlankMsg)
	}
	res, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.GetProductReq{
		ProductId: req.ID,
		UserId:    int32(userID),
	})
	if err != nil {
		l.Logger.Errorf("call rpc ProductRpc.GetProduct failed", logx.Field("err", err))
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}
	if res.StatusCode != code.Success {
		// 提示用户
		return nil, errors.New(int(res.StatusCode), res.StatusMsg)
	}

	resp = &types.GetProductByIDResp{
		ID:          int64(res.Product.Id),
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Picture:     res.Product.Picture,
		Stock:       res.Product.Stock,
		Price:       res.Product.Price,
		Sold:        res.Product.Sold,
		Categories:  res.Product.Categories,
	}

	return resp, err
}
