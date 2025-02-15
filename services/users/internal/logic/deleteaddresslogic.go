package logic

import (
	"context"
	"database/sql"

	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddressLogic {
	return &DeleteAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除用户地址
func (l *DeleteAddressLogic) DeleteAddress(in *users.DeleteAddressRequest) (*users.DeleteAddressResponse, error) {
	// todo: add your logic here and delete this line

	_, err := l.svcCtx.AddressModel.FindOne(l.ctx, in.AddressId)
	if err != nil {
		if err == sql.ErrNoRows {
			l.Logger.Infow(code.UserAddressNotFoundMsg, logx.Field("address_id", in.AddressId))
			return &users.DeleteAddressResponse{
				StatusCode: code.UserAddressNotFound,
				StatusMsg:  code.UserAddressNotFoundMsg,
			}, nil
		}
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("err", err))
		return &users.DeleteAddressResponse{
			StatusCode: code.ServerError,
			StatusMsg:  code.ServerErrorMsg,
		}, err
	}

	err = l.svcCtx.AddressModel.DeleteByAddressIdandUserId(l.ctx, in.AddressId, in.UserId)
	if err != nil {
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("err", err))
		return &users.DeleteAddressResponse{
			StatusCode: code.ServerError,
			StatusMsg:  code.ServerErrorMsg,
		}, err
	}

	return &users.DeleteAddressResponse{
		StatusCode: code.DeleteUserAddressSuccess,
		StatusMsg:  code.DeleteUserAddressSuccessMsg,
	}, nil
}
