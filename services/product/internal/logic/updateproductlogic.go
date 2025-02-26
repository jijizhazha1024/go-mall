package logic

import (
	"context"
	"database/sql"
	"fmt"
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
	// 1. 第一次删除缓存
	cacheKey := fmt.Sprintf("product:%d", in.Product.Id)
	if _, err := l.svcCtx.RedisClient.Del(cacheKey); err != nil {
		l.Logger.Errorw("product delete cache failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Product.Id))
		return nil, err
	}
	pictureUrl, err := UploadImage(in.Product.Picture, l.svcCtx.Config)
	if err != nil {
		l.Logger.Errorw("product picture upload failed",
			logx.Field("err", err))
		return nil, err
	}
	productRes := &product2.Products{
		Id:          int64(in.Product.Id),
		Name:        in.Product.Name,
		Description: sql.NullString{String: in.Product.Description, Valid: in.Product.Description != ""},
		Picture:     sql.NullString{String: pictureUrl, Valid: pictureUrl != ""},
		Price:       in.Product.Price,
		UpdatedAt:   time.Time{},
	}
	// 2. 使用 Transact 开启事务
	if err := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		// 得到图片对应url
		// 3. 更新商品记录：通过 withSession 生成支持事务的 updateModel 实例
		updateModel := product2.NewProductsModel(l.svcCtx.Mysql).WithSession(session)
		err = updateModel.Update(l.ctx, productRes)
		if err != nil {
			return fmt.Errorf("更新商品失败: %v", err)
		}

		// 4. 删除全部商品id的关联信息：生成基于事务的 product_categoriesModel 实例
		productCategoriesmodel := product_categories.NewProductCategoriesModel(l.svcCtx.Mysql).WithSession(session)
		if err := productCategoriesmodel.DeleteByProductId(l.ctx, int64(in.Product.Id)); err != nil {
			return fmt.Errorf("删除商品分类关联信息失败: %v", err)
		}

		// 5. 重新添加商品分类关联信息
		for _, category_id := range in.Product.Categories {
			categoryId, err := strconv.ParseInt(category_id, 10, 64)
			if err != nil {
				return fmt.Errorf("解析分类 ID 失败: %v", err)
			}
			if _, err := productCategoriesmodel.Insert(l.ctx, &product_categories.ProductCategories{
				ProductId:  sql.NullInt64{Int64: int64(in.Product.Id), Valid: int64(in.Product.Id) != 0},
				CategoryId: sql.NullInt64{Int64: categoryId, Valid: categoryId != 0},
			}); err != nil {
				return fmt.Errorf("插入商品分类关联信息失败: %v", err)
			}
		}

		return nil
	}); err != nil {
		l.Logger.Errorw("product update failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Product.Id))
		return &product.UpdateProductResp{
			StatusCode: uint32(code.ProductUpdateFailed),
			StatusMsg:  code.ProductUpdateFailedMsg,
		}, nil
	}
	// 7. 更新Elasticsearch记录
	docID := fmt.Sprintf("%d", in.Product.Id)

	// 构建更新请求（自动处理doc包装）
	updateResp, err := l.svcCtx.EsClient.Update().
		Index(biz.ProductEsIndexName).
		Id(docID).
		Doc(productRes).
		Refresh("true").
		DocAsUpsert(true). // 如果文档不存在则创建
		Do(l.ctx)          // 使用服务上下文
	if err != nil {
		// 处理文档不存在的情况（404错误）
		if elastic.IsNotFound(err) {
			l.Logger.Infow("尝试更新不存在的ES文档",
				logx.Field("product_id", in.Product.Id),
				logx.Field("doc_id", docID))
		} else {
			l.Logger.Errorw("product es update failed",
				logx.Field("err", err),
				logx.Field("product_id", in.Product.Id))
			return nil, err
		}
	} else {
		l.Logger.Infow("ES文档更新成功",
			logx.Field("index", updateResp.Index),
			logx.Field("doc_id", updateResp.Id),
			logx.Field("version", updateResp.Version))
	}
	// 9. 返回成功响应
	return &product.UpdateProductResp{
		StatusCode: uint32(code.ProductUpdated),
		StatusMsg:  code.ProductUpdatedMsg,
		I:          int64(in.Product.Id),
	}, nil
}
