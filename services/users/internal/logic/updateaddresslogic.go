package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/dal/model/user_address"
	"jijizhazha1024/go-mall/services/audit/audit"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

	//判断address——id和user——id是否存在

	_, err := l.svcCtx.AddressModel.GetUserAddressExistsByIdAndUserId(l.ctx, in.AddressId, int32(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			l.Logger.Infow("update address not found", logx.Field("address_id", in.AddressId), logx.Field("user_id", in.UserId))
			return &users.UpdateAddressResponse{
				StatusMsg:  code.UserAddressNotFoundMsg,
				StatusCode: code.UserAddressNotFound,
			}, nil
		}
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("user_id", in.UserId), logx.Field("err", err))
		return &users.UpdateAddressResponse{}, err
	}

	//判断修改后的地址是否是默认地址

	if in.IsDefault {

		addresses, err := l.svcCtx.AddressModel.FindAllByUserId(l.ctx, int32(in.UserId))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				l.Logger.Infow("update address is default, but user has no address", logx.Field("user_id", in.UserId))
				return &users.UpdateAddressResponse{
					StatusMsg:  code.UserAddressNotFoundMsg,
					StatusCode: code.UserAddressNotFound,
				}, nil
			}
			l.Logger.Errorw(code.ServerErrorMsg, logx.Field("user_id", in.UserId), logx.Field("err", err))
			return &users.UpdateAddressResponse{}, err
		}
		// 将所有地址的IsDefault字段设置为false+
		if err := l.svcCtx.Model.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {

			err := l.svcCtx.AddressModel.BatchUpdateDeFaultWithSession(ctx, session, addresses)
			if err != nil {
				return err
			}
			_, err = l.svcCtx.AddressModel.UpdateWithSession(l.ctx, session, &user_address.UserAddresses{

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
				UserId:          int64(in.UserId),
			})
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			l.Logger.Errorw("update address is__default is false, but update address failed", logx.Field("address_id", in.AddressId), logx.Field("err", err))
			return &users.UpdateAddressResponse{}, err
		}
	} else {
		err = l.svcCtx.AddressModel.Update(l.ctx, &user_address.UserAddresses{

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
			UserId:          int64(in.UserId),
		})
		if err != nil {
			l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("err", err))
			return &users.UpdateAddressResponse{}, err
		}

	}

	addressData, err := l.svcCtx.AddressModel.FindOne(l.ctx, int64(in.AddressId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			l.Logger.Infow("update address is not default, but address not found", logx.Field("address_id", in.AddressId), logx.Field("err", err))
			return &users.UpdateAddressResponse{
				StatusMsg:  code.UserAddressNotFoundMsg,
				StatusCode: code.UserAddressNotFound,
			}, nil
		}
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("err", err))
		return &users.UpdateAddressResponse{}, err
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
	newDataBytes, err := json.Marshal(data)
	if err != nil {
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("err", err))
		return &users.UpdateAddressResponse{}, err
	}

	newData := string(newDataBytes)

	//审计操作
	_, err = l.svcCtx.AuditRpc.CreateAuditLog(l.ctx, &audit.CreateAuditLogReq{

		UserId:            uint32(in.UserId),
		ActionType:        biz.Update,
		TargetTable:       "user_addresses",
		ActionDescription: "用户地址更新",
		ServiceName:       "users",
		TargetId:          int64(in.AddressId),
		ClientIp:          "127.0.0.1",
		NewData:           newData,
	})
	if err != nil {
		l.Logger.Infow(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("err", err))

	}

	return &users.UpdateAddressResponse{

		Data: data,
	}, nil
}
