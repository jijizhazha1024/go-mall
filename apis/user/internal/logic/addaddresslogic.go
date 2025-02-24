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

type AddAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAddressLogic {
	return &AddAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAddressLogic) AddAddress(req *types.AddAddressRequest) (resp *types.AddAddressResponse, err error) {

	//校验
	if req.City == "" || req.DetailedAddress == "" || req.PhoneNumber == "" || req.Province == "" {

		l.Logger.Errorw("用户信息为空", logx.Field("err", err))
		return nil, errors.New(code.Fail, "user informaition empty")

	}

	user_id := l.ctx.Value(biz.UserIDKey).(uint32)
	addaddressresp, err := l.svcCtx.UserRpc.AddAddress(l.ctx, &users.AddAddressRequest{

		UserId:          user_id,
		RecipientName:   req.RecipientName,
		Province:        req.Province,
		City:            req.City,
		PhoneNumber:     req.PhoneNumber,
		DetailedAddress: req.DetailedAddress,

		IsDefault: req.IsDefault,
	})

	if err != nil {

		l.Logger.Errorw("call rpc add address add failed", logx.Field("err", err))

		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	} else if addaddressresp.StatusMsg != "" {

		return nil, errors.New(int(addaddressresp.StatusCode), addaddressresp.StatusMsg)

	}

	Addressid := types.AddressData{
		AddressID:       addaddressresp.Data.AddressId,
		RecipientName:   addaddressresp.Data.RecipientName,
		PhoneNumber:     addaddressresp.Data.PhoneNumber,
		Province:        addaddressresp.Data.Province,
		City:            addaddressresp.Data.City,
		DetailedAddress: addaddressresp.Data.DetailedAddress,
		IsDefault:       addaddressresp.Data.IsDefault,
	}

	resp = &types.AddAddressResponse{
		Data: Addressid,
	}

	return
}
