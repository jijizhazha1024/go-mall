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

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户方法
func (l *UpdateUserLogic) UpdateUser(in *users.UpdateUserRequest) (*users.UpdateUserResponse, error) {
	// todo: add your logic here and delete this line

	update_user, err := l.svcCtx.UsersModel.FindOne(l.ctx, int64(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logx.Infow("update user not found", logx.Field("err", err),
				logx.Field("user_id", in.UserId))
			return &users.UpdateUserResponse{
				StatusCode: code.UserNotFound,
				StatusMsg:  code.UserNotFoundMsg,
			}, nil

		}
		logx.Errorw(code.ServerErrorMsg, logx.Field("err", err), logx.Field("user_id", in.UserId))
		return &users.UpdateUserResponse{}, nil

	}

	if update_user.UserDeleted {

		logx.Infow(" update user have deleted", logx.Field("user_id", in.UserId), logx.Field("user_id", in.UserId))

		return &users.UpdateUserResponse{
			StatusCode: code.UserHaveDeleted,
			StatusMsg:  code.UserHaveDeletedMsg,
		}, nil
	}

	err = l.svcCtx.UsersModel.UpdateUserName(l.ctx, int64(in.UserId), in.UsrName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logx.Infow("upate user not found", logx.Field("err", err),
				logx.Field("user_id", in.UserId))

			return &users.UpdateUserResponse{
				StatusCode: code.UserNotFound,
				StatusMsg:  code.UserNotFoundMsg,
			}, nil

		}

		return &users.UpdateUserResponse{}, err

	}

	//审计操作
	_, err = l.svcCtx.AuditRpc.CreateAuditLog(l.ctx, &audit.CreateAuditLogReq{

		UserId:            uint32(in.UserId),
		ActionType:        biz.Update,
		TargetTable:       "user",
		ActionDescription: "用户更新",
		ServiceName:       "users",
		OldData:           update_user.Username.String,
		NewData:           in.UsrName,
	})
	if err != nil {
		logx.Infow("create audit log failed", logx.Field("err", err))
		return &users.UpdateUserResponse{
			StatusCode: code.AuditUpdateuserFailed,
			StatusMsg:  code.AuditUpdateuserFailedMsg,
		}, nil

	}

	return &users.UpdateUserResponse{

		UserId: in.UserId,

		UserName: in.UsrName,
	}, nil

}
