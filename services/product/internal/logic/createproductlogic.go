package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"io/ioutil"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	pc "jijizhazha1024/go-mall/dal/model/products/product_categories"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"
	"strconv"
	"strings"
	"time"

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
	// todo: add your logic here and delete this line
	// 1. 敏感词校验
	if err := l.checkSensitiveWords(in.Name); err != nil {
		l.Logger.Errorw("product sensitive word failed",
			logx.Field("err", err))
		return &product.CreateProductResp{
			StatusCode: uint32(code.ProductSensitiveWordFailed),
			StatusMsg:  code.ProductSensitiveWordFailedMsg,
		}, err
	}
	var product_Id int64
	var picture_url string
	// 2. 使用 Transact 开启事务
	err := l.svcCtx.Mysql.Transact(func(session sqlx.Session) error {
		// 通过 withSession 生成支持事务的 productModel 实例
		productModel := product2.NewProductsModel(l.svcCtx.Mysql).WithSession(session)
		// 通过 withSession 生成支持事务的 productCategoriesModel 实例
		productCategoriesModel := pc.NewProductCategoriesModel(l.svcCtx.Mysql).WithSession(session)
		// 得到图片对应url
		picture_url, err := UploadImage(in.Picture, l.svcCtx.Config)
		if err != nil {
			return err
		}
		// 创建 Products 结构体实例
		productRes := &product2.Products{
			Name:        in.Name,
			Description: sql.NullString{String: in.Description, Valid: in.Description != ""},
			Picture:     sql.NullString{String: picture_url, Valid: picture_url != ""},
			Price:       float64(in.Price), // 注意类型转换
			Stock:       in.Stock,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		result, err := productModel.Insert(l.ctx, productRes)
		if err != nil {
			return fmt.Errorf("插入商品主表失败: %v", err)
		}
		product_Id, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("获取商品 ID 失败: %v", err)
		}

		// 3. 插入商品分类关联信息
		for _, category_id := range in.Categories {
			categoryId, err := strconv.ParseInt(category_id, 10, 64)
			if err != nil {
				return fmt.Errorf("解析分类 ID 失败: %v", err)
			}

			p := &pc.ProductCategories{
				ProductId:  sql.NullInt64{Int64: product_Id, Valid: product_Id != 0},
				CategoryId: sql.NullInt64{Int64: categoryId, Valid: categoryId != 0},
			}
			if _, err := productCategoriesModel.Insert(l.ctx, p); err != nil {
				return fmt.Errorf("插入商品分类关联信息失败: %v", err)
			}
		}

		return nil
	})

	// 4. 处理事务错误
	if err != nil {
		l.Logger.Errorw("product creation failed",
			logx.Field("err", err))
		return &product.CreateProductResp{
			StatusCode: uint32(code.ProductCreationFailed),
			StatusMsg:  code.ProductCreationFailedMsg,
		}, err
	}

	// 5. 将商品信息添加到 ES
	esDoc := map[string]interface{}{
		"id":          product_Id,
		"name":        in.Name,
		"description": in.Description,
		"picture":     picture_url,
		"price":       in.Price,
		"categories":  in.Categories,
	}
	var newBody string
	if newBody, err = mustJSON(esDoc); err != nil {
		l.Logger.Errorw("mustJSON err",
			logx.Field("err", err))
		return &product.CreateProductResp{
			StatusCode: uint32(code.ProductCreationFailed),
			StatusMsg:  code.ProductCreationFailedMsg,
		}, err
	}
	req := esapi.IndexRequest{
		Index:      biz.ProductEsIndexName, // ES 索引名，
		DocumentID: fmt.Sprintf("%d", product_Id),
		Body:       strings.NewReader(newBody),
		Refresh:    "true",
	}

	res, err := req.Do(l.ctx, l.svcCtx.Es)
	if err != nil {
		l.Logger.Errorw("product es creation failed",
			logx.Field("err", err))
		return &product.CreateProductResp{
			StatusCode: uint32(code.EsFailed),
			StatusMsg:  code.EsFailedMag,
		}, err
	}
	defer res.Body.Close()
	// 检查响应是否包含错误
	if res.IsError() {
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			l.Logger.Errorf("读取 Elasticsearch 响应体失败: %v", readErr)
		} else {
			l.Logger.Errorf("创建 Elasticsearch 记录时返回错误响应: %s", string(body))
		}
		return &product.CreateProductResp{
			StatusCode: uint32(code.EsFailed),
			StatusMsg:  code.EsFailedMag,
		}, err

	}
	return &product.CreateProductResp{
		StatusCode: uint32(code.ProductCreated),
		StatusMsg:  code.ProductCreatedMsg,
		I:          product_Id,
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
