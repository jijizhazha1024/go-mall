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

type UpdateAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAddressLogic {
	return &UpdateAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAddressLogic) UpdateAddress(req *types.UpdateAddressRequest) (resp *types.UpdateAddressResponse, err error) {

	if req.City == ""|| req.DetailedAddress == "" || req.PhoneNumber == "" || req.Province == "" {

		l.Logger.Errorf("用户信息为空", logx.Field("err", err))
		return nil, errors.New(code.Fail, "user informaition empty")

	}

	user_id := l.ctx.Value(biz.UserIDKey).(uint32)
	updateAddressresp, err := l.svcCtx.UserRpc.UpdateAddress(l.ctx, &users.UpdateAddressRequest{
		RecipientName:   req.RecipientName,
		PhoneNumber:     req.PhoneNumber,
		Province:        req.Province,
		City:            req.City,
		DetailedAddress: req.DetailedAddress,
		IsDefault:       req.IsDefault,
		AddressId:       req.AddressID,
		UserId:          user_id,
	})
	if err != nil {

		l.Logger.Errorf("call rpc updateaddress failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	} else {
		if updateAddressresp.StatusCode != code.UpdateUserAddressSuccess {
			l.Logger.Errorf("updeteaddress rpc failed", logx.Field("status_code", updateAddressresp.StatusCode), logx.Field("status_msg", updateAddressresp.StatusMsg))
			return nil, errors.New(int(updateAddressresp.StatusCode), updateAddressresp.StatusMsg)
		}
	}

	resp = &types.UpdateAddressResponse{
		Data: types.AddressData{
			AddressID:       updateAddressresp.Data.AddressId,
			RecipientName:   updateAddressresp.Data.RecipientName,
			PhoneNumber:     updateAddressresp.Data.PhoneNumber,
			Province:        updateAddressresp.Data.Province,
			City:            updateAddressresp.Data.City,
			DetailedAddress: updateAddressresp.Data.DetailedAddress,
			IsDefault:       updateAddressresp.Data.IsDefault,
			CreatedAt:       updateAddressresp.Data.CreatedAt,
			UpdatedAt:       updateAddressresp.Data.UpdatedAt,
		},
	}

	return resp, nil
}
