package logic

import (
	"context"
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
	// todo: add your logic here and delete this line
	product_id := in.Id
	productModel := product2.NewProductsModel(l.svcCtx.Mysql)
	exist, err := productModel.FindProductIsExist(l.ctx, product_id)
	if err != nil {
		l.Logger.Errorw("Failed to select data",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return nil, err
	}
	return &product.IsExistProductResp{
		StatusCode: uint32(code.ProductInfoRetrieved),
		StatusMsg:  code.ProductInfoRetrievedMsg,
		Exist:      exist,
	}, nil
}
