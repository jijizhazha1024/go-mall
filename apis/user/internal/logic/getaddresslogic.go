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

type GetAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressLogic {
	return &GetAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAddressLogic) GetAddress(req *types.GetAddressRequest) (resp *types.GetAddressResponse, err error) {

	user_id := l.ctx.Value(biz.UserIDKey).(uint32)
	getaddressresp, err := l.svcCtx.UserRpc.GetAddress(l.ctx, &users.GetAddressRequest{
		UserId:    user_id,
		AddressId: req.AddressID,
	})

	if err != nil {
		l.Logger.Errorf("调用 rpc 获取地址失败", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	} else {
		if getaddressresp.StatusCode != code.GetUserAddressSuccess {
			l.Logger.Errorf("调用 rpc 获取地址失败", logx.Field("status_code", getaddressresp.StatusCode), logx.Field("status_msg", getaddressresp.StatusMsg))
			return nil, errors.New(int(getaddressresp.StatusCode), getaddressresp.StatusMsg)
		}
	}

	// 创建响应对象并填充数据
	resp = &types.GetAddressResponse{
		Data: types.AddressData{
			AddressID:       getaddressresp.Data.AddressId,
			RecipientName:   getaddressresp.Data.RecipientName,
			PhoneNumber:     getaddressresp.Data.PhoneNumber,
			Province:        getaddressresp.Data.Province,
			City:            getaddressresp.Data.City,
			DetailedAddress: getaddressresp.Data.DetailedAddress,
			IsDefault:       getaddressresp.Data.IsDefault,
		},
	}

	return resp, nil

}
