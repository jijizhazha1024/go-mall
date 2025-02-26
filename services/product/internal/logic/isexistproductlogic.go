package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"

	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsExistProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsExistProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsExistProductLogic {
	return &IsExistProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 判断商品是否存在
func (l *IsExistProductLogic) IsExistProduct(in *product.IsExistProductReq) (*product.IsExistProductResp, error) {
	productModel := product2.NewProductsModel(l.svcCtx.Mysql)
	exist, err := productModel.FindProductIsExist(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			// 不存在并不属于错误，所以这里不需要返回错误，由调用端返回信息
			return &product.IsExistProductResp{
				StatusCode: uint32(code.ProductNotFoundInventory),
				StatusMsg:  code.ProductNotFoundInventoryMsg,
			}, nil
		}
		l.Logger.Errorw("Failed to select data",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return nil, err
	}
	return &product.IsExistProductResp{
		Exist: exist,
	}, nil
}
