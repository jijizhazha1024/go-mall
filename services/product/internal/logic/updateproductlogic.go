package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	"jijizhazha1024/go-mall/dal/model/products/product_categories"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"
	"strconv"
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

// UpdateProduct 修改商品
func (l *UpdateProductLogic) UpdateProduct(in *product.UpdateProductReq) (*product.UpdateProductResp, error) {
	// 1. 第一次删除缓存
	cacheKey := fmt.Sprintf(biz.ProductIDKey, in.Id)
	if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
		l.Logger.Errorw("product delete cache failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return nil, err
	}
	var pictureUrl string
	if len(in.Picture) != 0 {
		zone, _ := storage.GetZone(l.svcCtx.Config.QiNiu.AccessKey, l.svcCtx.Config.QiNiu.Bucket)
		url, err := UploadImage(in.Picture, zone, l.svcCtx.Config)
		if err != nil {
			l.Logger.Errorw("product picture upload failed",
				logx.Field("err", err))
			return nil, err
		}
		pictureUrl = url
	}

	productRes := &product2.Products{
		Id:          in.Id,
		Name:        in.Name,
		Description: sql.NullString{String: in.Description, Valid: in.Description != ""},
		Picture:     sql.NullString{String: pictureUrl, Valid: pictureUrl != ""},
		Price:       in.Price,
	}
	res := &product.UpdateProductResp{}
	// 2. 使用 Transact 开启事务
	if err := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		updateModel := product2.NewProductsModel(l.svcCtx.Mysql).WithSession(session)
		if err := updateModel.Update(l.ctx, productRes); err != nil {
			return err
		}
		exist, err := updateModel.FindProductIsExist(l.ctx, in.Id)
		if err != nil {
			return err
		}
		if !exist {
			res.StatusCode = code.ProductNotFound
			res.StatusMsg = code.ProductNotFoundMsg
			return nil
		}
		// 4. 删除全部商品id的关联信息：生成基于事务的 product_categoriesModel 实例
		productCategoriesmodel := product_categories.NewProductCategoriesModel(l.svcCtx.Mysql).WithSession(session)
		if err := productCategoriesmodel.DeleteByProductId(l.ctx, in.Id); err != nil {
			return err
		}

		// 5. 重新添加商品分类关联信息
		for _, categoryId := range in.Categories {
			categoryId, err := strconv.ParseInt(categoryId, 10, 64)
			if err != nil {
				return err
			}
			if _, err := productCategoriesmodel.Insert(l.ctx, &product_categories.ProductCategories{
				ProductId:  sql.NullInt64{Int64: in.Id, Valid: int64(in.Id) != 0},
				CategoryId: sql.NullInt64{Int64: categoryId, Valid: categoryId != 0},
			}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		l.Logger.Errorw("product update failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return nil, err
	}

	// 3. 更新Elasticsearch记录
	if _, err := l.svcCtx.EsClient.Update().
		Index(biz.ProductEsIndexName).
		Id(strconv.Itoa(int(in.Id))).
		Doc(productRes).
		Refresh("true").
		DocAsUpsert(true). // 如果文档不存在则创建
		Do(l.ctx); err != nil && !elastic.IsNotFound(err) {
		l.Logger.Errorw("product es update failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return res, nil
	}
	res.Id = in.Id
	return res, nil
}
