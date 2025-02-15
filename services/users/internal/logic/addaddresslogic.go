package logic

import (
	"context"
	"database/sql"

	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/dal/model/user_address"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAddressLogic {
	return &AddAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAddressLogic) AddAddress(in *users.AddAddressRequest) (*users.AddAddressResponse, error) {
	// todo: add your logic here and delete this line
	phonenumber := sql.NullString{
		String: in.PhoneNumber,
		Valid:  in.Province != "",
	}
	province := sql.NullString{
		String: in.Province,
		Valid:  in.Province != "",
	}

	result, err := l.svcCtx.AddressModel.Insert(l.ctx, &user_address.UserAddresses{
		UserId:          int64(in.UserId),
		RecipientName:   in.RecipientName,
		PhoneNumber:     phonenumber,
		Province:        province,
		City:            in.City,
		DetailedAddress: in.DetailedAddress,
		IsDefault:       in.IsDefault,
	})

	if err != nil {
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("err", err))
		return &users.AddAddressResponse{
			StatusMsg:  code.AddUserAddressFailedMsg,
			StatusCode: code.AddUserAddressFailed,
		}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("err", err))
		return &users.AddAddressResponse{
			StatusMsg:  code.AddUserAddressFailedMsg,
			StatusCode: code.AddUserAddressFailed,
		}, err
	}
	data := &users.AddressData{
		AddressId:       int32(id),
		RecipientName:   in.RecipientName,
		PhoneNumber:     in.PhoneNumber,
		Province:        in.Province,
		City:            in.City,
		DetailedAddress: in.DetailedAddress,
		IsDefault:       in.IsDefault,
	}
	return &users.AddAddressResponse{
		StatusMsg:  code.AddUserAddressSuccessMsg,
		StatusCode: code.AddUserAddressSuccess,
		Data:       data,
	}, nil

}
