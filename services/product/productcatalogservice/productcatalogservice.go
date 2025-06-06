// Code generated by goctl. DO NOT EDIT.
// Source: product.proto

package productcatalogservice

import (
	"context"

	"jijizhazha1024/go-mall/services/product/product"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateProductReq              = product.CreateProductReq
	CreateProductResp             = product.CreateProductResp
	DeleteProductReq              = product.DeleteProductReq
	DeleteProductResp             = product.DeleteProductResp
	GetAllProductsReq             = product.GetAllProductsReq
	GetAllProductsResp            = product.GetAllProductsResp
	GetProductReq                 = product.GetProductReq
	GetProductResp                = product.GetProductResp
	IsExistProductReq             = product.IsExistProductReq
	IsExistProductResp            = product.IsExistProductResp
	Product                       = product.Product
	QueryProductReq               = product.QueryProductReq
	QueryProductReq_Paginator     = product.QueryProductReq_Paginator
	QueryProductReq_Price         = product.QueryProductReq_Price
	RecommendProductReq           = product.RecommendProductReq
	RecommendProductReq_Paginator = product.RecommendProductReq_Paginator
	UpdateProductReq              = product.UpdateProductReq
	UpdateProductResp             = product.UpdateProductResp

	ProductCatalogService interface {
		// 根据商品id得到商品详细信息
		GetProduct(ctx context.Context, in *GetProductReq, opts ...grpc.CallOption) (*GetProductResp, error)
		// 添加新商品
		CreateProduct(ctx context.Context, in *CreateProductReq, opts ...grpc.CallOption) (*CreateProductResp, error)
		// 修改商品
		UpdateProduct(ctx context.Context, in *UpdateProductReq, opts ...grpc.CallOption) (*UpdateProductResp, error)
		// 删除商品
		DeleteProduct(ctx context.Context, in *DeleteProductReq, opts ...grpc.CallOption) (*DeleteProductResp, error)
		// 分页得到全部商品
		GetAllProduct(ctx context.Context, in *GetAllProductsReq, opts ...grpc.CallOption) (*GetAllProductsResp, error)
		// 判断商品是否存在
		IsExistProduct(ctx context.Context, in *IsExistProductReq, opts ...grpc.CallOption) (*IsExistProductResp, error)
		// 根据条件查询商品
		QueryProduct(ctx context.Context, in *QueryProductReq, opts ...grpc.CallOption) (*GetAllProductsResp, error)
		RecommendProduct(ctx context.Context, in *RecommendProductReq, opts ...grpc.CallOption) (*GetAllProductsResp, error)
	}

	defaultProductCatalogService struct {
		cli zrpc.Client
	}
)

func NewProductCatalogService(cli zrpc.Client) ProductCatalogService {
	return &defaultProductCatalogService{
		cli: cli,
	}
}

// 根据商品id得到商品详细信息
func (m *defaultProductCatalogService) GetProduct(ctx context.Context, in *GetProductReq, opts ...grpc.CallOption) (*GetProductResp, error) {
	client := product.NewProductCatalogServiceClient(m.cli.Conn())
	return client.GetProduct(ctx, in, opts...)
}

// 添加新商品
func (m *defaultProductCatalogService) CreateProduct(ctx context.Context, in *CreateProductReq, opts ...grpc.CallOption) (*CreateProductResp, error) {
	client := product.NewProductCatalogServiceClient(m.cli.Conn())
	return client.CreateProduct(ctx, in, opts...)
}

// 修改商品
func (m *defaultProductCatalogService) UpdateProduct(ctx context.Context, in *UpdateProductReq, opts ...grpc.CallOption) (*UpdateProductResp, error) {
	client := product.NewProductCatalogServiceClient(m.cli.Conn())
	return client.UpdateProduct(ctx, in, opts...)
}

// 删除商品
func (m *defaultProductCatalogService) DeleteProduct(ctx context.Context, in *DeleteProductReq, opts ...grpc.CallOption) (*DeleteProductResp, error) {
	client := product.NewProductCatalogServiceClient(m.cli.Conn())
	return client.DeleteProduct(ctx, in, opts...)
}

// 分页得到全部商品
func (m *defaultProductCatalogService) GetAllProduct(ctx context.Context, in *GetAllProductsReq, opts ...grpc.CallOption) (*GetAllProductsResp, error) {
	client := product.NewProductCatalogServiceClient(m.cli.Conn())
	return client.GetAllProduct(ctx, in, opts...)
}

// 判断商品是否存在
func (m *defaultProductCatalogService) IsExistProduct(ctx context.Context, in *IsExistProductReq, opts ...grpc.CallOption) (*IsExistProductResp, error) {
	client := product.NewProductCatalogServiceClient(m.cli.Conn())
	return client.IsExistProduct(ctx, in, opts...)
}

// 根据条件查询商品
func (m *defaultProductCatalogService) QueryProduct(ctx context.Context, in *QueryProductReq, opts ...grpc.CallOption) (*GetAllProductsResp, error) {
	client := product.NewProductCatalogServiceClient(m.cli.Conn())
	return client.QueryProduct(ctx, in, opts...)
}

func (m *defaultProductCatalogService) RecommendProduct(ctx context.Context, in *RecommendProductReq, opts ...grpc.CallOption) (*GetAllProductsResp, error) {
	client := product.NewProductCatalogServiceClient(m.cli.Conn())
	return client.RecommendProduct(ctx, in, opts...)
}
