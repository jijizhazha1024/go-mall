package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	"jijizhazha1024/go-mall/dal/model/products/product_categories"
	"strconv"
	"strings"
	"time"

	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改商品
func (l *UpdateProductLogic) UpdateProduct(in *product.UpdateProductReq) (*product.UpdateProductResp, error) {
	var picture_url string
	// 1. 第一次删除缓存
	cacheKey := fmt.Sprintf("product:%d", in.Product.Id)
	if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
		l.Logger.Errorw("product delete cache failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Product.Id))
		return &product.UpdateProductResp{
			StatusCode: uint32(code.ProductCacheFailed),
			StatusMsg:  code.ProductCacheFailedMsg,
		}, err
	}

	// 2. 使用 Transact 开启事务
	err := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		// 得到图片对应url
		picture_url, err := UploadImage(in.Product.Picture, l.svcCtx.Config)
		if err != nil {
			return err
		}
		// 3. 更新商品记录：通过 withSession 生成支持事务的 updateModel 实例
		updateModel := product2.NewProductsModel(l.svcCtx.Mysql).WithSession(session)
		err = updateModel.Update(l.ctx, &product2.Products{
			Id:          int64(in.Product.Id),
			Name:        in.Product.Name,
			Description: sql.NullString{String: in.Product.Description, Valid: in.Product.Description != ""},
			Picture:     sql.NullString{String: picture_url, Valid: picture_url != ""},
			Price:       float64(in.Product.Price),
			Stock:       in.Product.Stock,
			UpdatedAt:   time.Time{},
		})
		if err != nil {
			return fmt.Errorf("更新商品失败: %v", err)
		}

		// 4. 删除全部商品id的关联信息：生成基于事务的 product_categoriesModel 实例
		product_categoriesModel := product_categories.NewProductCategoriesModel(l.svcCtx.Mysql).WithSession(session)
		if err := product_categoriesModel.DeleteByProductId(l.ctx, int64(in.Product.Id)); err != nil {
			return fmt.Errorf("删除商品分类关联信息失败: %v", err)
		}

		// 5. 重新添加商品分类关联信息
		for _, category_id := range in.Product.Categories {
			categoryId, err := strconv.ParseInt(category_id, 10, 64)
			if err != nil {
				return fmt.Errorf("解析分类 ID 失败: %v", err)
			}
			if _, err := product_categoriesModel.Insert(l.ctx, &product_categories.ProductCategories{
				ProductId:  sql.NullInt64{Int64: int64(in.Product.Id), Valid: int64(in.Product.Id) != 0},
				CategoryId: sql.NullInt64{Int64: categoryId, Valid: categoryId != 0},
			}); err != nil {
				return fmt.Errorf("插入商品分类关联信息失败: %v", err)
			}
		}

		return nil
	})

	// 6. 处理事务错误
	if err != nil {
		l.Logger.Errorw("product update failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Product.Id))
		return &product.UpdateProductResp{
			StatusCode: uint32(code.ProductUpdateFailed),
			StatusMsg:  code.ProductUpdateFailedMsg,
		}, err
	}
	// 7. 更新Elasticsearch记录
	// Elasticsearch索引名称
	indexName := biz.ProductEsIndexName
	// Elasticsearch文档ID
	docID := fmt.Sprintf("%d", in.Product.Id)
	// 构造Elasticsearch文档
	esDoc := map[string]interface{}{
		"id":          in.Product.Id,
		"name":        in.Product.Name,
		"description": in.Product.Description,
		"picture":     picture_url,
		"price":       in.Product.Price,
		"categories":  in.Product.Categories,
	}
	// 构造正确的更新请求体
	updateBody := map[string]interface{}{
		"doc": esDoc,
	}
	var ubstring string
	if ubstring, err = mustJSON(updateBody); err != nil {
		l.Logger.Errorw("mustJSON err",
			logx.Field("err", err))
		return &product.UpdateProductResp{
			StatusCode: uint32(code.ProductUpdateFailed),
			StatusMsg:  code.ProductUpdateFailedMsg,
		}, err
	}
	// 创建Elasticsearch更新请求
	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       strings.NewReader(ubstring),
		Refresh:    "true",
	}

	// 发送请求
	res, err := req.Do(context.Background(), l.svcCtx.Es)
	if err != nil {
		l.Logger.Errorw("product es update failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Product.Id))
		return &product.UpdateProductResp{
			StatusCode: uint32(code.EsFailed),
			StatusMsg:  code.EsFailedMag,
		}, err
	}
	defer res.Body.Close()

	// 检查响应
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			l.Logger.Errorf("解析Elasticsearch响应失败: %v", err)
		}
		l.Logger.Errorw("product es update body failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Product.Id))
		return &product.UpdateProductResp{
			StatusCode: uint32(code.EsFailed),
			StatusMsg:  code.EsFailedMag,
		}, err
	}

	// 8. 延迟第二次删除缓存
	go func() {
		time.Sleep(500 * time.Millisecond)
		if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
			l.Logger.Errorf("第二次删除缓存失败: %v", err)
		}
	}()
	// 9. 返回成功响应
	return &product.UpdateProductResp{
		StatusCode: uint32(code.ProductUpdated),
		StatusMsg:  code.ProductUpdatedMsg,
		I:          int64(in.Product.Id),
	}, nil
}
