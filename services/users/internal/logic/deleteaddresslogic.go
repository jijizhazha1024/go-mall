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

	address1, err := l.svcCtx.AddressModel.GetUserAddressExistsByIdAndUserId(l.ctx, in.AddressId, int32(in.UserId))
	if err != nil {

		l.Logger.Errorw("find address by id and user id failed", logx.Field("address_id", in.AddressId), logx.Field("user_id", in.UserId), logx.Field("err", err))
		return nil, err
	}
	if !address1 {
		return &users.DeleteAddressResponse{
			StatusMsg:  code.UserAddressNotFoundMsg,
			StatusCode: code.UserAddressNotFound,
		}, nil

	}

	err = l.svcCtx.AddressModel.DeleteByAddressIdandUserId(l.ctx, in.AddressId, int32(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return &users.DeleteAddressResponse{

				StatusCode: code.UserAddressNotFound,
				StatusMsg:  code.UserAddressNotFoundMsg,
			}, nil
		}
		l.Logger.Errorw("delete address failed", logx.Field("address_id", in.AddressId), logx.Field("user_id", in.UserId), logx.Field("err", err))
		return nil, err
	}
	//添加审计服务
	auditreq := audit.CreateAuditLogReq{
		UserId:            uint32(in.UserId),
		ActionType:        biz.Delete,
		TargetTable:       "user_address",
		ActionDescription: "删除用户地址",
		ClientIp:          in.Ip,
		TargetId:          int64(in.AddressId),
		ServiceName:       "users",
	}
	_, err = l.svcCtx.AuditRpc.CreateAuditLog(l.ctx, &auditreq)
	if err != nil {
		l.Logger.Infow("add address audit failed", logx.Field("err", err),
			logx.Field("body", auditreq))

	}

	return nil, nil
}
