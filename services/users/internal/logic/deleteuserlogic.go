package logic

import (
	"context"
	"database/sql"
	"errors"

	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/internal/users_biz"
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

			return users_biz.HandleDeleteUsererror(code.UserNotFoundMsg, code.UserNotFound, nil)
		}
		logx.Errorw(code.ServerErrorMsg, logx.Field("err", err), logx.Field("user_id", in.UserId))
		return users_biz.HandleDeleteUsererror(code.ServerErrorMsg, code.ServerError, err)
	}
	// 删除用户
	if exituser.UserDeleted {
		l.Logger.Infow("delete user have deleted", logx.Field("user_id", in.UserId))

		return users_biz.HandleDeleteUsererror(code.UserHaveDeletedMsg, code.UserHaveDeleted, nil)
	}
	err = l.svcCtx.UsersModel.UpdateDeletebyId(l.ctx, int64(in.UserId), true)
	if err != nil {
		l.Logger.Infow("delete update deletebyid failed", logx.Field("err", err),
			logx.Field("user_id", in.UserId))

		return users_biz.HandleDeleteUsererror(code.UserDeletionFailedMsg, code.UserDeletionFailed, nil)
	}

	return users_biz.HandleDeleteUserResp(code.UserDeletedMsg, code.UserDeleted, in.UserId)
}
