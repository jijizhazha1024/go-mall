package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	pc "jijizhazha1024/go-mall/dal/model/products/product_categories"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 添加新商品
func (l *CreateProductLogic) CreateProduct(in *product.CreateProductReq) (*product.CreateProductResp, error) {
	// 1. 敏感词校验
	if err := l.checkSensitiveWords(in.Name); err != nil {
		l.Logger.Errorw("product sensitive word failed",
			logx.Field("err", err))
		return nil, err
	}
	var productId int64
	pictureUrl, err := UploadImage(in.Picture, l.svcCtx.Config)
	if err != nil {
		l.Logger.Errorw("product picture upload failed",
			logx.Field("err", err))
		return nil, err
	}
	// 创建 Products 结构体实例
	productRes := &product2.Products{
		Name:        in.Name,
		Description: sql.NullString{String: in.Description, Valid: in.Description != ""},
		Picture:     sql.NullString{String: pictureUrl, Valid: pictureUrl != ""},
		Price:       in.Price, // 注意类型转换
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// 2. 使用 Transact 开启事务
	if err := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		// 通过 withSession 生成支持事务的 productModel 实例
		productModel := product2.NewProductsModel(l.svcCtx.Mysql).WithSession(session)
		// 通过 withSession 生成支持事务的 productCategoriesModel 实例
		productCategoriesModel := pc.NewProductCategoriesModel(l.svcCtx.Mysql).WithSession(session)
		// 得到图片对应url
		result, err := productModel.Insert(l.ctx, productRes)
		if err != nil {
			return fmt.Errorf("插入商品主表失败: %v", err)
		}
		productId, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("获取商品 ID 失败: %v", err)
		}
		// 3. 插入商品分类关联信息
		for _, categoryId := range in.Categories {
			categoryId, err := strconv.ParseInt(categoryId, 10, 64)
			if err != nil {
				return fmt.Errorf("解析分类 ID 失败: %v", err)
			}

			p := &pc.ProductCategories{
				ProductId:  sql.NullInt64{Int64: productId, Valid: productId != 0},
				CategoryId: sql.NullInt64{Int64: categoryId, Valid: categoryId != 0},
			}
			if _, err := productCategoriesModel.Insert(l.ctx, p); err != nil {
				return fmt.Errorf("插入商品分类关联信息失败: %v", err)
			}
		}
		return nil
	}); err != nil {
		l.Logger.Errorw("product creation failed",
			logx.Field("err", err))
		return &product.CreateProductResp{
			StatusCode: uint32(code.ProductCreationFailed),
			StatusMsg:  code.ProductCreationFailedMsg,
		}, nil
	}
	// 创建文档（自动JSON序列化）
	if _, err = l.svcCtx.EsClient.Index().
		Index(biz.ProductEsIndexName).
		Id(fmt.Sprintf("%d", productId)).
		BodyJson(map[string]interface{}{
			"name":        productRes.Name,
			"description": productRes.Description.String,
			"picture":     productRes.Picture.String,
			"price":       productRes.Price,
			"categories":  in.Categories,
			"created_at":  productRes.CreatedAt,
			"updated_at":  productRes.UpdatedAt,
		}).
		Refresh("true").
		Do(l.ctx); err != nil {
		l.Logger.Errorw("product es creation failed",
			logx.Field("err", err))
		return nil, err
	}

	return &product.CreateProductResp{
		I: productId,
	}, nil
}

func (l *CreateProductLogic) checkSensitiveWords(text string) error {
	// 敏感词过滤逻辑
	// 目前仅使用简单的字符串匹配

	if text == "敏感词" {
		return fmt.Errorf("包含敏感词")
	}
	return nil
}

// mustJSON 辅助函数，用于将结构体转换为 JSON 字符串
func mustJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
