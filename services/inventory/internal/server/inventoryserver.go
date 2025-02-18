// Code generated by goctl. DO NOT EDIT.
// Source: inventory.proto

package server

import (
	"context"

	"jijizhazha1024/go-mall/services/inventory/internal/logic"
	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"
)

type InventoryServer struct {
	svcCtx *svc.ServiceContext
	inventory.UnimplementedInventoryServer
}

func NewInventoryServer(svcCtx *svc.ServiceContext) *InventoryServer {
	return &InventoryServer{
		svcCtx: svcCtx,
	}
}

// GetInventory 查询库存
func (s *InventoryServer) GetInventory(ctx context.Context, in *inventory.GetInventoryReq) (*inventory.GetInventoryResp, error) {
	l := logic.NewGetInventoryLogic(ctx, s.svcCtx)
	return l.GetInventory(in)
}

// UpdateInventory 增加库存
func (s *InventoryServer) UpdateInventory(ctx context.Context, in *inventory.InventoryReq) (*inventory.InventoryResp, error) {
	l := logic.NewUpdateInventoryLogic(ctx, s.svcCtx)
	return l.UpdateInventory(in)
}

// DecreaseInventory 扣减库存
func (s *InventoryServer) DecreaseInventory(ctx context.Context, in *inventory.InventoryReq) (*inventory.InventoryResp, error) {
	l := logic.NewDecreaseInventoryLogic(ctx, s.svcCtx)
	return l.DecreaseInventory(in)
}

// ReturnInventory 归还库存
func (s *InventoryServer) ReturnInventory(ctx context.Context, in *inventory.InventoryReq) (*inventory.InventoryResp, error) {
	l := logic.NewReturnInventoryLogic(ctx, s.svcCtx)
	return l.ReturnInventory(in)
}
