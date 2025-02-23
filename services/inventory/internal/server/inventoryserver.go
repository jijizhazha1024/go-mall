// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5
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

// GetInventory 查询库存，缓存不在，再去数据库查
func (s *InventoryServer) GetInventory(ctx context.Context, in *inventory.GetInventoryReq) (*inventory.GetInventoryResp, error) {
	l := logic.NewGetInventoryLogic(ctx, s.svcCtx)
	return l.GetInventory(in)
}

func (s *InventoryServer) GetPreInventory(ctx context.Context, in *inventory.GetPreInventoryReq) (*inventory.GetPreInventoryResp, error) {
	l := logic.NewGetPreInventoryLogic(ctx, s.svcCtx)
	return l.GetPreInventory(in)
}

// UpdateInventory 增加库存，修改库存数量（直接修改）
func (s *InventoryServer) UpdateInventory(ctx context.Context, in *inventory.InventoryReq) (*inventory.InventoryResp, error) {
	l := logic.NewUpdateInventoryLogic(ctx, s.svcCtx)
	return l.UpdateInventory(in)
}

// DecreaseInventory 预扣减库存，此时并非真实扣除库存，而是在缓存进行--操作
func (s *InventoryServer) DecreasePreInventory(ctx context.Context, in *inventory.InventoryReq) (*inventory.InventoryResp, error) {
	l := logic.NewDecreasePreInventoryLogic(ctx, s.svcCtx)
	return l.DecreasePreInventory(in)
}

// DecreaseInventory 真实扣减库存（支付成功时）
func (s *InventoryServer) DecreaseInventory(ctx context.Context, in *inventory.InventoryReq) (*inventory.InventoryResp, error) {
	l := logic.NewDecreaseInventoryLogic(ctx, s.svcCtx)
	return l.DecreaseInventory(in)
}

// ReturnPreInventory 退还预扣减的库存（）
func (s *InventoryServer) ReturnPreInventory(ctx context.Context, in *inventory.InventoryReq) (*inventory.InventoryResp, error) {
	l := logic.NewReturnPreInventoryLogic(ctx, s.svcCtx)
	return l.ReturnPreInventory(in)
}

// ReturnInventory 退还库存（支付失败时）
func (s *InventoryServer) ReturnInventory(ctx context.Context, in *inventory.InventoryReq) (*inventory.InventoryResp, error) {
	l := logic.NewReturnInventoryLogic(ctx, s.svcCtx)
	return l.ReturnInventory(in)
}
