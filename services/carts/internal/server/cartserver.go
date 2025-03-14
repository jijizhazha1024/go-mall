// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5
// Source: carts.proto

package server

import (
	"context"

	"jijizhazha1024/go-mall/services/carts/carts"
	"jijizhazha1024/go-mall/services/carts/internal/logic"
	"jijizhazha1024/go-mall/services/carts/internal/svc"
)

type CartServer struct {
	svcCtx *svc.ServiceContext
	carts.UnimplementedCartServer
}

func NewCartServer(svcCtx *svc.ServiceContext) *CartServer {
	return &CartServer{
		svcCtx: svcCtx,
	}
}

func (s *CartServer) CartItemList(ctx context.Context, in *carts.UserInfo) (*carts.CartItemListResponse, error) {
	l := logic.NewCartItemListLogic(ctx, s.svcCtx)
	return l.CartItemList(in)
}

func (s *CartServer) CreateCartItem(ctx context.Context, in *carts.CartItemRequest) (*carts.CreateCartResponse, error) {
	l := logic.NewCreateCartItemLogic(ctx, s.svcCtx)
	return l.CreateCartItem(in)
}
func (s *CartServer) SubCartItem(ctx context.Context, in *carts.CartItemRequest) (*carts.SubCartResponse, error) {
	l := logic.NewSubCartItemLogic(ctx, s.svcCtx)
	return l.SubCartItem(in)
}

func (s *CartServer) DeleteCartItem(ctx context.Context, in *carts.CartItemRequest) (*carts.EmptyCartResponse, error) {
	l := logic.NewDeleteCartItemLogic(ctx, s.svcCtx)
	return l.DeleteCartItem(in)
}
