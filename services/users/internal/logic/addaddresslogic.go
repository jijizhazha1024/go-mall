package logic

import (
	"context"
	"database/sql"
	"errors"

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
	var addresses []*user_address.UserAddresses

	if in.IsDefault {
		var err error
		addresses, err = l.svcCtx.AddressModel.FindAllByUserId(l.ctx, int32(in.UserId))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				l.Logger.Infow("update address is default, but user has no address", logx.Field("user_id", in.UserId))
				return &users.AddAddressResponse{
					StatusMsg:  code.UserAddressNotFoundMsg,
					StatusCode: code.UserAddressNotFound,
				}, nil
			}
			l.Logger.Errorw(code.ServerErrorMsg, logx.Field("user_id", in.UserId), logx.Field("err", err))
			return &users.AddAddressResponse{
				StatusMsg:  code.ServerErrorMsg,
				StatusCode: code.ServerError,
			}, err
		}
	}

	if addresses != nil {
		//如果存在，则将其设为非默认地址然后增加新的默认地址
		var result sql.Result

		if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {

			err := l.svcCtx.AddressModel.BatchUpdateDeFaultWithSession(ctx, session, addresses)
			if err != nil {
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
			return &users.AddAddressResponse{}, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			l.Logger.Errorw("id insert failed", logx.Field("user_id", in.UserId), logx.Field("err", err))
			return &users.AddAddressResponse{}, err
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

		//添加审计服务

		return &users.AddAddressResponse{

			Data: data,
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
			return &users.AddAddressResponse{}, err
		}

		// 获取插入ID并赋值给外部变量
		id, err := result.LastInsertId()
		if err != nil {
			l.Logger.Errorw("id insert failed", logx.Field("user_id", in.UserId), logx.Field("err", err))
			return &users.AddAddressResponse{}, err
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
		//添加审计服务
		return &users.AddAddressResponse{

			Data: data,
		}, nil

	}

}
