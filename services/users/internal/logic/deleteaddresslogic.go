package logic

import (
	"context"
	"database/sql"
	"errors"

	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/audit/audit"
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

	_, err := l.svcCtx.AddressModel.GetUserAddressExistsByIdAndUserId(l.ctx, in.AddressId, int32(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
		if errors.Is(err, sql.ErrNoRows) {
			l.Logger.Infow("deleteaddress address not found", logx.Field("err", err), logx.Field("address_id", in.AddressId), logx.Field("user_id", in.UserId))
			return &users.DeleteAddressResponse{

				StatusCode: code.UserAddressNotFound,
				StatusMsg:  code.UserAddressNotFoundMsg,
			}, nil
		}
		l.Logger.Errorw(code.ServerErrorMsg, logx.Field("address_id", in.AddressId), logx.Field("user_id", in.UserId), logx.Field("err", err))
		return &users.DeleteAddressResponse{}, err
	}
	//添加审计服务
	_, err = l.svcCtx.AuditRpc.CreateAuditLog(l.ctx, &audit.CreateAuditLogReq{

		UserId:            uint32(in.UserId),
		ActionType:        biz.Delete,
		TargetTable:       "user",
		ActionDescription: "删除用户地址",
		TargetId:          int64(in.AddressId),
		ServiceName:       "users",
	})
	if err != nil {
		l.Logger.Infow("add address audit failed", logx.Field("err", err),
			logx.Field("user_id", in.UserId))
		return &users.DeleteAddressResponse{
			StatusCode: code.AuditDeleteaddressFailed,
			StatusMsg:  code.AuditDeleteaddressFailedMsg,
		}, nil
	}

	return &users.DeleteAddressResponse{}, nil
}
