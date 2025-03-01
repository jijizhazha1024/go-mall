package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/common/consts/biz"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"
	"strconv"
)

func ScanHotProducts(svcCtx *svc.ServiceContext, ctx context.Context) (err error) {
	logc := logx.WithContext(ctx)
	products, err := svcCtx.RedisClient.Zrevrange(biz.ProductRedisPVName, 0, 99)
	if err != nil {
		logc.Infow("Failed to get top 100 hot products",
			logx.Field("err", err),
		)
		return err
	}
	productModel := product2.NewProductsModel(svcCtx.Mysql)
	// 打印或处理获取到的商品 ID
	for _, productID := range products {
		productId, err := strconv.ParseInt(productID, 10, 64)
		if err != nil {
			logc.Infow("Failed to convert product id to int",
				logx.Field("productId", productID),
				logx.Field("err", err),
			)
			continue
		}
		productData, err := productModel.FindOne(ctx, productId)
		if err != nil {
			logc.Infow("Failed to find hot product",
				logx.Field("productId", productID),
				logx.Field("err", err),
			)
			continue
		}
		// 构造响应
		resp := &product.Product{
			Id:          uint32(productData.Id),
			Name:        productData.Name,
			Description: productData.Description.String,
			Picture:     productData.Picture.String,
			Price:       productData.Price,
			Categories:  nil,
		}

		// 将数据缓存到Redis中
		data, err := json.Marshal(resp)
		cacheData := string(data)
		if err != nil {
			logc.Infow("Failed to unmarshal data",
				logx.Field("err", err),
			)
			return err
		}
		if err = svcCtx.RedisClient.Set(productID, cacheData); err != nil {
			logc.Infow("Failed to set hot product cache",
				logx.Field("productId", productID),
				logx.Field("err", err),
			)
		}
	}
	return nil

}
