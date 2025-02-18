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
	//判断address——id和user——id是否存在

	_, err := l.svcCtx.AddressModel.GetUserAddressbyIdAndUserId(l.ctx, in.AddressId, int32(in.UserId))
	if err != nil {
		if err == sql.ErrNoRows {
			l.Logger.Infow("delete address not found", logx.Field("address_id", in.AddressId), logx.Field("user_id", in.UserId))
			return &users.DeleteAddressResponse{
				StatusMsg:  code.UserAddressNotFoundMsg,
				StatusCode: code.UserAddressNotFound,
			}, nil
		}
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("user_id", in.UserId), logx.Field("err", err))
		return &users.DeleteAddressResponse{}, err
	}

	err = l.svcCtx.AddressModel.DeleteByAddressIdandUserId(l.ctx, in.AddressId, int32(in.UserId))
	if err != nil {
		if err == sql.ErrNoRows {
			l.Logger.Infow("deleteaddress address not found", logx.Field("err", err), logx.Field("address_id", in.AddressId), logx.Field("user_id", in.UserId))
			return &users.DeleteAddressResponse{

				StatusCode: code.UserAddressNotFound,
				StatusMsg:  code.UserAddressNotFoundMsg,
			}, nil
		}
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("user_id", in.UserId), logx.Field("err", err))
		return &users.DeleteAddressResponse{}, err
	}

	return &users.DeleteAddressResponse{}, nil
}
