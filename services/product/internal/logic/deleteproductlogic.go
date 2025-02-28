package logic

import (
	"context"
	"github.com/olivere/elastic/v7"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	"jijizhazha1024/go-mall/dal/model/products/product_categories"
	"strconv"
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
	// 删除商品对应pv
	// 商品 ID
	//ProductIDKey
	cacheKey := strconv.Itoa(int(in.Id))
	// 删除哈希表中的商品 ID 字段
	_, err := l.svcCtx.RedisClient.Zrem(biz.ProductRedisPVName, cacheKey)
	if err != nil {
		l.Logger.Errorw("从 Redis 哈希表中删除商品失败",
			logx.Field("productId", in.Id),
			logx.Field("err", err))
		return nil, err
	}
	// 1. 第一次删除缓存
	if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
		l.Logger.Errorw("product delete cache failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return nil, err
	}
	res := &product.DeleteProductResp{}
	// 2. 使用 Transact 开启事务
	if err = l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		// 3. 删除商品记录：通过 withSession 生成支持事务的 deleteModel 实例
		deleteModel := product2.NewProductsModel(l.svcCtx.Mysql).WithSession(session)
		exist, err := deleteModel.FindProductIsExist(l.ctx, in.Id)
		if err != nil {
			return err
		}
		if !exist {
			res.StatusCode = code.ProductNotFound
			res.StatusMsg = code.ProductNotFoundMsg
			return nil
		}
		if err := deleteModel.Delete(l.ctx, in.Id); err != nil {
			return err
		}
		// 4. 删除商品分类关系：同样生成基于事务的 deleteCategoryModel 实例
		deleteCategoryModel := product_categories.NewProductCategoriesModel(l.svcCtx.Mysql).WithSession(session)
		if err := deleteCategoryModel.DeleteByProductId(l.ctx, in.Id); err != nil {
			return err
		}
		return nil
	}); err != nil {
		l.Logger.Errorw("product delete  failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return nil, err
	}
	if res.StatusCode != code.Success {
		return res, nil
	}

	// 6. 删除es记录
	// 构建删除请求
	if _, err := l.svcCtx.EsClient.Delete().
		Index(biz.ProductEsIndexName).
		Id(strconv.Itoa(int(in.Id))).
		Refresh("true").
		Do(l.ctx); err != nil && elastic.IsNotFound(err) {
		l.Logger.Errorw("product es delete failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return nil, err
	}

	// 7. 延迟第二次删除缓存
	go func() {
		time.Sleep(500 * time.Millisecond)
		if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
			l.Logger.Errorf("第二次删除缓存失败: %v", err)
		}
	}()
	// 8. 返回成功响应
	return res, nil
}
