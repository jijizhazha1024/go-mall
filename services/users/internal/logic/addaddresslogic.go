package logic

import (
	"context"
	"database/sql"

	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/dal/model/user_address"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

	phonenumber := sql.NullString{
		String: in.PhoneNumber,
		Valid:  in.PhoneNumber != "",
	}
	province := sql.NullString{
		String: in.Province,
		Valid:  in.Province != "",
	}

	//查找用户是否已经存在默认地址
	defaultAddress, err := l.svcCtx.AddressModel.FindDefaultByUserId(l.ctx, int32(in.UserId))
	if err != nil {
		l.Logger.Infow("find default address failed", logx.Field("user_id", in.UserId), logx.Field("err", err))

	}

	if defaultAddress != nil {
		//如果存在，则将其设为非默认地址然后增加新的默认地址
		var result sql.Result

		defaultAddress.IsDefault = false

		if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {

			_, err = l.svcCtx.AddressModel.UpdateWithSession(l.ctx, session, defaultAddress)
			if err != nil {
				l.Logger.Errorw("update default address failed", logx.Field("user_id", in.UserId), logx.Field("err", err))
				return err
			}

			result, err = l.svcCtx.AddressModel.InsertWithSession(l.ctx, session, &user_address.UserAddresses{

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
				return err
			}

			return nil
		}); err != nil {
			l.Logger.Errorw("update and add default address failed", logx.Field("user_id", in.UserId), logx.Field("err", err))
			return &users.AddAddressResponse{
				StatusMsg:  code.AddUserAddressFailedMsg,
				StatusCode: code.AddUserAddressFailed,
			}, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			l.Logger.Errorw("id insert failed", logx.Field("user_id", in.UserId), logx.Field("err", err))
			return &users.AddAddressResponse{
				StatusMsg:  code.AddUserAddressFailedMsg,
				StatusCode: code.AddUserAddressFailed,
			}, err
		}

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

	} else { //不存在 增加新的默认地址
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
		if err != nil {
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

}
