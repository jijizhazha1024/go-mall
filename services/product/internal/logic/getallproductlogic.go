package logic

import (
	"context"
	product2 "jijizhazha1024/go-mall/dal/model/products/product"

	"jijizhazha1024/go-mall/services/product/internal/svc"
	"jijizhazha1024/go-mall/services/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllProductLogic {
	return &GetAllProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页得到全部商品
func (l *GetAllProductLogic) GetAllProduct(in *product.GetAllProductsReq) (*product.GetAllProductsResp, error) {
	page := in.Page
	pageSize := in.PageSize
	// 计算偏移量
	offset := (page - 1) * pageSize
	productModel := product2.NewProductsModel(l.svcCtx.Mysql)
	// 查询商品列表
	products, err := productModel.FindPage(l.ctx, int(offset), int(pageSize))
	if err != nil {
		l.Logger.Errorw("product select  failed",
			logx.Field("err", err))
		return nil, err
	}

	// 查询总记录数
	total, err := productModel.Count(l.ctx)
	if err != nil {
		l.Logger.Errorw("product select  failed",
			logx.Field("err", err))
		return nil, err
	}

	// 转换为响应结构
	var result []*product.Product
	for _, p := range products {
		result = append(result, &product.Product{
			Id:    uint32(p.Id),
			Name:  p.Name,
			Price: float32(p.Price),
			Stock: p.Stock,
		})
	}

	// 构造响应对象
	resp := &product.GetAllProductsResp{
		Products: result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	return resp, nil
}
