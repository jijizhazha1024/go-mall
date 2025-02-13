package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"
)

type GetProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据商品id得到商品详细信息
func (l *GetProductLogic) GetProduct(in *product.GetProductReq) (*product.GetProductResp, error) {
	// todo: add your logic here and delete this line
	product_id := in.Id
	productModel := product2.NewProductsModel(l.svcCtx.Mysql)
	// 在redis中维护商品的访问频率次数 PV
	// 检查商品 ID 是否存在
	redisKey := biz.ProductRedisPVName
	cacheKey := fmt.Sprintf("%d", product_id)
	_, err := l.svcCtx.RedisClient.Zincrby(redisKey, 1, cacheKey)
	if err != nil {
		l.Logger.Errorw("自增商品的访问次数失败",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return &product.GetProductResp{
			StatusCode: uint32(code.ProductCacheFailed),
			StatusMsg:  code.ProductCacheFailedMsg,
		}, err
	}
	// 从Redis中获取数据
	cacheData, err := l.svcCtx.RedisClient.Get(cacheKey)
	if err != nil {
		l.Logger.Errorw("get product from cache failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return &product.GetProductResp{
			StatusCode: uint32(code.ProductCacheFailed),
			StatusMsg:  code.ProductCacheFailedMsg,
		}, err
	}

	// 如果Redis中有数据，直接反序列化并返回
	if cacheData != "" {
		var productRes product.Product
		if err := json.Unmarshal([]byte(cacheData), &productRes); err != nil {
			l.Logger.Errorw("Failed to unmarshal data",
				logx.Field("err", err),
				logx.Field("product_id", in.Id))
			return &product.GetProductResp{
				StatusCode: uint32(code.ProductCacheFailed),
				StatusMsg:  code.ProductCacheFailedMsg,
			}, err
		}
		return &product.GetProductResp{
			StatusCode: uint32(code.ProductInfoRetrieved),
			StatusMsg:  code.ProductInfoRetrievedMsg,
			Product:    &productRes,
		}, nil
	}

	// 如果Redis中没有数据，从数据库中获取
	productData, err := productModel.FindOne(l.ctx, int64(product_id))
	if err != nil {
		l.Logger.Errorw("Failed to find product from database",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return &product.GetProductResp{
			StatusCode: uint32(code.ProductInfoRetrievalFailed),
			StatusMsg:  code.ProductInfoRetrievalFailedMsg,
		}, err
	}
	// 构造响应
	resp := &product.Product{
		Id:          uint32(productData.Id),
		Name:        productData.Name,
		Description: productData.Description.String,
		Picture:     productData.Picture.String,
		Price:       float32(productData.Price),
		Stock:       productData.Stock,
		Categories:  nil,
	}

	// 将数据缓存到Redis中
	data, err := json.Marshal(resp)
	cacheData = string(data)
	if err != nil {
		l.Logger.Errorw("Failed to unmarshal data",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return &product.GetProductResp{
			StatusCode: uint32(code.ProductCacheFailed),
			StatusMsg:  code.ProductCacheFailedMsg,
		}, err
	}
	err = l.svcCtx.RedisClient.Set(cacheKey, cacheData)
	if err != nil {
		l.Logger.Errorw("Failed to save redis data",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return &product.GetProductResp{
			StatusCode: uint32(code.ProductCacheFailed),
			StatusMsg:  code.ProductCacheFailedMsg,
		}, err
	}
	return &product.GetProductResp{
		StatusCode: uint32(code.ProductInfoRetrieved),
		StatusMsg:  code.ProductInfoRetrievedMsg,
		Product:    resp,
	}, nil
}
