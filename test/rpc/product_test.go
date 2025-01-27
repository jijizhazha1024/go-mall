package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/product/product"
	"testing"
)

var product_client product.ProductCatalogServiceClient

func initproduct() {
	conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%d", biz.ProductRpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	product_client = product.NewProductCatalogServiceClient(conn)
}

func TestProductsRpc(t *testing.T) {
	initproduct()
	resp, err := product_client.CreateProduct(context.Background(), &product.CreateProductReq{
		Name:        "小米",
		Description: "手机信息",
		Picture:     "..",
		Price:       123.23,
		Stock:       5000,
		Categories:  []string{"1", "2", "3"},
	})
	if err != nil {
		t.Fatal(err)

	}
	fmt.Println(" success", resp)
	t.Log(" success", resp)
}
