package logic

import (
	"context"

	"jijizhazha1024/go-mall/apis/user/internal/svc"
	"jijizhazha1024/go-mall/apis/user/internal/types"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
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
	// 调用用户 RPC 获取地址列表
	user_id := l.ctx.Value(biz.UserIDKey).(uint32)
	listaddressresp, err := l.svcCtx.UserRpc.ListAddresses(l.ctx, &users.AllAddressLitstRequest{
		UserId: user_id,
	})

	if err != nil {
		l.Logger.Errorw("调用 rpc 获取地址列表失败", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	} else if listaddressresp.StatusMsg != "" {

		return nil, errors.New(int(listaddressresp.StatusCode), listaddressresp.StatusMsg)

	}

	// 创建响应对象并填充数据
	resp = &types.AddressListResponse{
		Data: make([]types.AddressData, 0),
	}
	for _, address := range listaddressresp.Data {
		resp.Data = append(resp.Data, types.AddressData{
			AddressID:       address.AddressId,
			RecipientName:   address.RecipientName,
			PhoneNumber:     address.PhoneNumber,
			Province:        address.Province,
			City:            address.City,
			DetailedAddress: address.DetailedAddress,
			IsDefault:       address.IsDefault,
			CreatedAt:       address.CreatedAt,
			UpdatedAt:       address.UpdatedAt,
		})
	}

	return resp, nil
}
