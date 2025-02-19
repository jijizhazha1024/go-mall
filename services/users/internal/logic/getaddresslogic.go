package logic

import (
	"context"
	"database/sql"
	"errors"

	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressLogic {
	return &GetAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取单个收货地址
func (l *GetAddressLogic) GetAddress(in *users.GetAddressRequest) (*users.GetAddressResponse, error) {
	// todo: add your logic here and delete this line

	address, err := l.svcCtx.AddressModel.GetUserAddressbyIdAndUserId(l.ctx, in.AddressId, int32(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			l.Logger.Infow("get address not found", logx.Field("user_id", in.UserId), logx.Field("address_id", in.AddressId))
			return &users.GetAddressResponse{
				StatusMsg:  code.UserAddressNotFoundMsg,
				StatusCode: code.UserAddressNotFound,
			}, nil
		}

		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("user_id", in.UserId), logx.Field("address_id", in.AddressId), logx.Field("err", err))
		return &users.GetAddressResponse{}, err
	}

	data := &users.AddressData{
		AddressId:       int32(address.AddressId),
		RecipientName:   address.RecipientName,
		PhoneNumber:     address.PhoneNumber.String,
		Province:        address.Province.String,
		City:            address.City,
		DetailedAddress: address.DetailedAddress,
		IsDefault:       address.IsDefault,
		CreatedAt:       address.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       address.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return &users.GetAddressResponse{

		Data: data,
	}, nil
}
