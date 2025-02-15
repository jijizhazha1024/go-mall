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

	//将当前用户的地址信息
	addresses, err := l.svcCtx.AddressModel.FindAllByUserId(l.ctx, in.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			l.Logger.Infow(code.UserAddressNotFoundMsg, logx.Field("user_id", in.UserId))
			return &users.UpdateAddressResponse{
				StatusMsg:  code.UserAddressNotFoundMsg,
				StatusCode: code.UserAddressNotFound,
			}, nil
		}
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("user_id", in.UserId), logx.Field("err", err))
		return &users.UpdateAddressResponse{
			StatusMsg:  code.ServerErrorMsg,
			StatusCode: code.ServerError,
		}, err
	}
	// 将所有地址的IsDefault字段设置为false+
	err = l.svcCtx.AddressModel.BatchUpdateDeFAULT(l.ctx, addresses)
	if err != nil {
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("user_id", in.UserId), logx.Field("err", err))
		return &users.UpdateAddressResponse{
			StatusMsg:  code.ServerErrorMsg,
			StatusCode: code.ServerError,
		}, err

	}

	_, err = l.svcCtx.AddressModel.Update(l.ctx, &user_address.UserAddresses{

		AddressId:     int64(in.AddressId),
		RecipientName: in.RecipientName,
		PhoneNumber: sql.NullString{
			String: string(in.PhoneNumber),
			Valid:  in.PhoneNumber != ""},
		Province: sql.NullString{
			String: string(in.Province),
			Valid:  in.Province != ""},
		City:            in.City,
		DetailedAddress: in.DetailedAddress,
		IsDefault:       in.IsDefault,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			l.Logger.Infow(code.UserAddressNotFoundMsg, logx.Field("address_id", in.AddressId), logx.Field("err", err))
			return &users.UpdateAddressResponse{
				StatusMsg:  code.UserAddressNotFoundMsg,
				StatusCode: code.UserAddressNotFound,
			}, nil
		}
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("err", err))
		return &users.UpdateAddressResponse{
			StatusMsg:  code.ServerErrorMsg,
			StatusCode: code.ServerError,
		}, nil
	}

	addressData, err := l.svcCtx.AddressModel.FindOne(l.ctx, in.AddressId)
	if err != nil {
		if err == sql.ErrNoRows {
			l.Logger.Infow(code.UserAddressNotFoundMsg, logx.Field("address_id", in.AddressId), logx.Field("err", err))
			return &users.UpdateAddressResponse{
				StatusMsg:  code.UserAddressNotFoundMsg,
				StatusCode: code.UserAddressNotFound,
			}, nil
		}
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("err", err))
		return &users.UpdateAddressResponse{
			StatusMsg:  code.ServerErrorMsg,
			StatusCode: code.ServerError,
		}, err
	}

	data := &users.AddressData{
		AddressId:       int32(addressData.AddressId),
		RecipientName:   addressData.RecipientName,
		PhoneNumber:     addressData.PhoneNumber.String,
		Province:        addressData.Province.String,
		City:            addressData.City,
		DetailedAddress: addressData.DetailedAddress,
		IsDefault:       addressData.IsDefault,
		CreatedAt:       addressData.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       addressData.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return &users.UpdateAddressResponse{
		StatusMsg:  code.UpdateUserAddressSuccessMsg,
		StatusCode: code.UpdateUserAddressSuccess,
		Data:       data,
	}, nil
}
