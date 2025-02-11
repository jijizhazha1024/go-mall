package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAddressLogic {
	return &UpdateAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改用户地址
func (l *UpdateAddressLogic) UpdateAddress(in *users.UpdateAddressRequest) (*users.UpdateAddressResponse, error) {
	// todo: add your logic here and delete this line

	return &users.UpdateAddressResponse{}, nil
}
