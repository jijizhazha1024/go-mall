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

type ListAddressesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAddressesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAddressesLogic {
	return &ListAddressesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取所有收货地址
func (l *ListAddressesLogic) ListAddresses(in *users.AllAddressLitstRequest) (*users.AddressListResponse, error) {

	alladdress, err := l.svcCtx.AddressModel.FindAllByUserId(l.ctx, int32(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			l.Logger.Infow("address list user address not found", logx.Field("user_id", in.UserId))
			return &users.AddressListResponse{
				StatusMsg:  code.UserAddressNotFoundMsg,
				StatusCode: code.UserAddressNotFound,
			}, nil
		}
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("user_id", in.UserId), logx.Field("err", err))
		return &users.AddressListResponse{}, err
	}
	addresslist := make([]*users.AddressData, 0)
	for _, address := range alladdress {
		addresslist = append(addresslist, &users.AddressData{
			AddressId:       int32(address.AddressId),
			DetailedAddress: address.DetailedAddress,
			City:            address.City,
			Province:        address.Province.String,
			IsDefault:       address.IsDefault,
			RecipientName:   address.RecipientName,
			PhoneNumber:     address.PhoneNumber.String,
			CreatedAt:       address.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       address.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &users.AddressListResponse{

		Data: addresslist,
	}, nil
}
