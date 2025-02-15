package logic

import (
	"context"

	"jijizhazha1024/go-mall/apis/user/internal/svc"
	"jijizhazha1024/go-mall/apis/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllAddressListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAllAddressListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllAddressListLogic {
	return &AllAddressListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AllAddressListLogic) AllAddressList(req *types.AllAddressListRequest) (resp *types.AddressListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
