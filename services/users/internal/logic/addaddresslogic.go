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
		Valid:  in.PhoneNumber != "",
	}
	province := sql.NullString{
		String: in.Province,
		Valid:  in.Province != "",
	}
	// 在外部作用域声明ID变量

	// 使用事务session执行插入
	result, err := l.svcCtx.AddressModel.Insert(l.ctx, &user_address.UserAddresses{
		UserId:          int64(in.UserId),
		DetailedAddress: in.DetailedAddress,
		City:            in.City,
		Province:        province,
		IsDefault:       in.IsDefault,
		RecipientName:   in.RecipientName,
		PhoneNumber:     phonenumber,
	})
	if err != nil {
		l.Logger.Errorw("add address failed", logx.Field("user_id", in.UserId), logx.Field("err", err))
		return &users.AddAddressResponse{
			StatusMsg:  code.AddUserAddressFailedMsg,
			StatusCode: code.AddUserAddressFailed,
		}, err
	}

	// 获取插入ID并赋值给外部变量
	id, err := result.LastInsertId()
	if err == nil {
		l.Logger.Errorw("id insert failed", logx.Field("user_id", in.UserId), logx.Field("err", err))
		return &users.AddAddressResponse{
			StatusMsg:  code.AddUserAddressFailedMsg,
			StatusCode: code.AddUserAddressFailed,
		}, err
	}

	// 构建返回数据（此时id已赋值）
	data := &users.AddressData{
		AddressId:       int32(id), // 使用事务中获取的ID
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
