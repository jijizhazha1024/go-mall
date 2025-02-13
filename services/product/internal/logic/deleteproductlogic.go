package logic

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"io/ioutil"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	"jijizhazha1024/go-mall/dal/model/products/product_categories"
	"time"

	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除商品
func (l *DeleteProductLogic) DeleteProduct(in *product.DeleteProductReq) (*product.DeleteProductResp, error) {
	// todo: add your logic here and delete this line
	// 删除商品对应pv
	redisKey := biz.ProductRedisPVName
	// 商品 ID
	productID := fmt.Sprintf("%d", in.Id)
	cacheKey := fmt.Sprintf("%d", in.Id)
	// 删除哈希表中的商品 ID 字段
	_, err := l.svcCtx.RedisClient.Zrem(redisKey, cacheKey)
	if err != nil {
		l.Logger.Errorw("从 Redis 哈希表中删除商品失败",
			logx.Field("productId", productID),
			logx.Field("err", err))
		return &product.DeleteProductResp{
			StatusCode: uint32(code.ProductPVCacheFailed),
			StatusMsg:  code.ProductPVCacheFailedMsg,
		}, err
	}
	// 1. 第一次删除缓存
	if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
		l.Logger.Errorw("product delete cache failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return &product.DeleteProductResp{
			StatusCode: uint32(code.ProductCacheFailed),
			StatusMsg:  code.ProductCacheFailedMsg,
		}, err
	}

	// 2. 使用 Transact 开启事务
	err = l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		// 3. 删除商品记录：通过 withSession 生成支持事务的 deleteModel 实例
		deleteModel := product2.NewProductsModel(l.svcCtx.Mysql).WithSession(session)
		if err := deleteModel.Delete(l.ctx, in.Id); err != nil {
			return fmt.Errorf("删除商品失败: %v", err)
		}
		// 4. 删除商品分类关系：同样生成基于事务的 deleteCategoryModel 实例
		deleteCategoryModel := product_categories.NewProductCategoriesModel((l.svcCtx.Mysql)).WithSession(session)
		if err := deleteCategoryModel.DeleteByProductId(l.ctx, in.Id); err != nil {
			return fmt.Errorf("删除商品分类关系失败: %v", err)
		}

		return nil
	})

	// 5. 处理事务错误
	if err != nil {
		l.Logger.Errorw("product delete  failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return &product.DeleteProductResp{
			StatusCode: uint32(code.ProductDeletionFailed),
			StatusMsg:  code.ProductDeletionFailedMsg,
		}, err
	}

	// 6. 删除es记录

	// Elasticsearch文档ID
	docID := fmt.Sprintf("%d", in.Id)

	// 删除Elasticsearch记录
	req := esapi.DeleteRequest{
		Index:      biz.ProductEsIndexName,
		DocumentID: docID,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), l.svcCtx.Es)
	if err != nil {
		l.Logger.Errorw("product es delete  failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return &product.DeleteProductResp{
			StatusCode: uint32(code.EsFailed),
			StatusMsg:  code.EsFailedMag,
		}, err
	}
	if res != nil {
		defer res.Body.Close()
	}
	// 检查响应是否包含错误
	if res.IsError() {
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			l.Logger.Errorf("读取 Elasticsearch 响应体失败: %v", readErr)
		} else {
			l.Logger.Errorf("删除 Elasticsearch 记录时返回错误响应: %s", string(body))
		}
		l.Logger.Errorw("product es resp read failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return &product.DeleteProductResp{
			StatusCode: uint32(code.EsFailed),
			StatusMsg:  code.EsFailedMag,
		}, err
	}

	// 7. 延迟第二次删除缓存
	go func() {
		time.Sleep(500 * time.Millisecond)
		if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
			l.Logger.Errorf("第二次删除缓存失败: %v", err)
		}
	}()
	// 8. 返回成功响应
	return &product.DeleteProductResp{
		StatusCode: uint32(code.ProductDeleted),
		StatusMsg:  code.ProductDeletedMsg,
	}, nil
}
