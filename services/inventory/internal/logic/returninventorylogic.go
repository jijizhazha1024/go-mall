package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReturnInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReturnInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReturnInventoryLogic {
	return &ReturnInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReturnInventory 退还库存（支付失败时）
func (l *ReturnInventoryLogic) ReturnInventory(in *inventory.InventoryReq) (*inventory.InventoryResp, error) {
	// todo: add your logic here and delete this line

	return &inventory.InventoryResp{}, nil
}
