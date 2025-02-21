package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/inventory/internal/svc"
	"jijizhazha1024/go-mall/services/inventory/inventory"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPreInventoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPreInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPreInventoryLogic {
	return &GetPreInventoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPreInventoryLogic) GetPreInventory(in *inventory.GetPreInventoryReq) (*inventory.GetPreInventoryResp, error) {
	//从缓存中查找库存

	return &inventory.GetPreInventoryResp{}, nil
}
