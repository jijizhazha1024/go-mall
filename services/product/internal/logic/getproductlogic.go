package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	gorse "jijizhazha1024/go-mall/common/utils/gorse"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"
	"jijizhazha1024/go-mall/services/inventory/inventory"
	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"
	"strconv"
	"time"
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

// GetProduct 根据商品id得到商品详细信息
func (l *GetProductLogic) GetProduct(in *product.GetProductReq) (*product.GetProductResp, error) {

	// 在redis中维护商品的访问频率次数 PV
	// 检查商品 ID 是否存在
	redisKey := biz.ProductRedisPVName
	cacheKey := fmt.Sprintf(biz.ProductIDKey, in.Id)
	_, err := l.svcCtx.RedisClient.Zincrby(redisKey, 1, cacheKey)
	if err != nil {
		// 这里可以只进行记录即可，可以无需返回,还是可以正常的进行执行的，不影响返回结果
		l.Logger.Errorw("自增商品的访问次数失败",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
	}
	// 从Redis中获取数据
	cacheData, err := l.svcCtx.RedisClient.Get(cacheKey)
	if err != nil {
		// 这里也是，可以想象一个场景，假设在请求redis时网络抖动了，导致请求失败，但是在后面还可以通过mysql获取数据
		l.Logger.Errorw("get product from cache failed",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
	}

	// 如果Redis中有数据且没有错误，直接反序列化并返回
	if err == nil && cacheData != "" {
		var productRes product.Product
		if err := json.Unmarshal([]byte(cacheData), &productRes); err == nil {
			// 序列化成功返回，查询库存，我们进行返回动态库存
			productRes.Stock, productRes.Sold = l.getRealTimeStockAndSold(int64(in.Id))
			return &product.GetProductResp{
				Product: &productRes,
			}, nil
		}
		// 序列失败 也是一样进行记录日志，因为在后面还可以从mysql查询，这样用用户体验感好点
		logx.Errorw("Failed to unmarshal data",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
	}

	// 如果Redis中没有数据，从数据库中获取

	productModel := product2.NewProductsModel(l.svcCtx.Mysql)
	productData, err := productModel.FindOne(l.ctx, int64(in.Id))
	// 存在错误直接返回，因为没有兜底的了。
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			// 不存在并不属于错误，所以这里不需要返回错误，由调用端返回信息
			return &product.GetProductResp{
				StatusCode: uint32(code.ProductNotFoundInventory),
				StatusMsg:  code.ProductNotFoundInventoryMsg,
			}, nil
		}
		l.Logger.Errorw("Failed to find product from database",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return nil, err
	}

	// --------------- 到这里就说明查询到部分数据了，进行组装操作 ---------------
	resp := &product.GetProductResp{
		Product: &product.Product{
			Id:          uint32(productData.Id),
			Name:        productData.Name,
			Description: productData.Description.String,
			Picture:     productData.Picture.String,
			Price:       productData.Price,
			CratedAt:    productData.CreatedAt.Format(time.DateTime),
			UpdatedAt:   productData.CreatedAt.Format(time.DateTime),
		},
	}

	// 在这里创建连接，懒惰创建连接。
	categories, err := l.svcCtx.CategoriesModel.FindCategoryNameByProductID(l.ctx, int64(in.Id))
	if err != nil {
		l.Logger.Errorw("Failed to find product_category from database",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		// 因为查询不完整，所以不需要写入缓存了，直接返回
		return resp, nil
	}

	resp.Product.Categories = categories
	// 到这里就说明数据是完整的，将数据缓存到Redis中
	data, err := json.Marshal(resp.Product)
	cacheData = string(data)
	if err != nil {
		l.Logger.Errorw("Failed to unmarshal data",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return resp, nil
	}
	// 设置合理的过期时间
	if err = l.svcCtx.RedisClient.SetexCtx(l.ctx, cacheKey, cacheData, biz.ProductIDKeyExpire); err != nil {
		l.Logger.Errorw("Failed to save redis data",
			logx.Field("err", err),
			logx.Field("product_id", in.Id))
		return resp, nil
	}

	// 查询实时库存
	resp.Product.Stock, resp.Product.Sold = l.getRealTimeStockAndSold(productData.Id)
	if in.UserId != 0 {
		go func() {
			// 插入反馈
			if _, err := l.svcCtx.GorseClient.InsertFeedback(l.ctx, []gorse.Feedback{
				{
					ItemId:       strconv.Itoa(int(productData.Id)),
					UserId:       strconv.Itoa(int(in.UserId)),
					Timestamp:    time.Now().Format(time.DateTime),
					FeedbackType: biz.ReadFeedBackType,
				},
			}); err != nil {
				l.Logger.Infow("Failed to insert feedback", logx.Field("err", err), logx.Field("product_id", productData.Id))
				return
			}
		}()
	}

	return resp, nil
}

// 抽取重复的库存查询逻辑
func (l *GetProductLogic) getRealTimeStockAndSold(productId int64) (int64, int64) {
	inventoryResp, err := l.svcCtx.InventoryRpc.GetInventory(l.ctx, &inventory.GetInventoryReq{
		ProductId: int32(productId),
	})
	if err != nil {
		l.Logger.Errorw("call rpc InventoryRpc.GetInventory failed", logx.Field("err", err), logx.Field("product_id", productId))
		return 0, 0 // 返回默认值或特殊标记
	}
	return inventoryResp.Inventory, inventoryResp.SoldCount
}
