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

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除用户方法
func (l *DeleteUserLogic) DeleteUser(in *users.DeleteUserRequest) (*users.DeleteUserResponse, error) {

	exituser, err := l.svcCtx.UsersModel.FindOne(l.ctx, int64(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			l.Logger.Infow("delete user not found", logx.Field("err", err),
				logx.Field("user_id", in.UserId))
			return &users.DeleteUserResponse{

				StatusCode: code.UserAddressNotFound,
				StatusMsg:  code.UserAddressNotFoundMsg,
			}, nil
		}
		logx.Errorw(code.ServerErrorMsg, logx.Field("err", err), logx.Field("user_id", in.UserId))
		return &users.DeleteUserResponse{}, err
	}
	// 删除用户
	if exituser.UserDeleted {
		l.Logger.Infow("delete user have deleted", logx.Field("user_id", in.UserId))
		return &users.DeleteUserResponse{

			StatusCode: code.UserHaveDeleted,
			StatusMsg:  code.UserHaveDeletedMsg,
		}, nil

	}
	err = l.svcCtx.UsersModel.UpdateDeletebyId(l.ctx, int64(in.UserId), true)
	if err != nil {
		l.Logger.Infow("delete update deletebyid failed", logx.Field("err", err),
			logx.Field("user_id", in.UserId))

		return &users.DeleteUserResponse{}, err
	}

	//添加审计服务
	_, err = l.svcCtx.AuditRpc.CreateAuditLog(l.ctx, &audit.CreateAuditLogReq{

		UserId:            uint32(in.UserId),
		ActionType:        biz.Delete,
		TargetTable:       "user",
		ActionDescription: "用户注销",
		TargetId:          int64(in.UserId),
		ServiceName:       "users",
		ClientIp:          "127.0.0.1",
	})
	if err != nil {
		l.Logger.Infow("add address audit failed", logx.Field("err", err),
			logx.Field("user_id", in.UserId))

	}
	return &users.DeleteUserResponse{
		UserId: in.UserId,
	}, nil
}
