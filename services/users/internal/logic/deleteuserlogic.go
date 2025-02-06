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
	// todo: add your logic here and delete this line

	// 查询用户是否存在
	exituser, err := l.svcCtx.UsersModel.FindOne(l.ctx, int64(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			l.Logger.Info(code.UserNotFoundMsg, in.UserId, err)
			return users_biz.HandleDeleteUsererror(code.UserNotFoundMsg, code.UserNotFound)
		}
		l.Logger.Error(code.ServerErrorMsg, err)
		return users_biz.HandleDeleteUsererror(code.ServerErrorMsg, code.ServerError)
	}
	// 删除用户
	if exituser.UserDeleted {
		l.Logger.Info(code.UserHaveDeletedMsg, in.UserId)
		return users_biz.HandleDeleteUsererror(code.UserHaveDeletedMsg, code.UserHaveDeleted)
	}
	err = l.svcCtx.UsersModel.UpdateDeletebyId(l.ctx, int64(in.UserId), true)
	if err != nil {
		l.Logger.Error(code.UserDeletionFailedMsg, err)
		return users_biz.HandleDeleteUsererror(code.UserDeletionFailedMsg, code.UserDeletionFailed)
	}

	return users_biz.HandleDeleteUserResp(code.UserDeletedMsg, code.UserDeleted, in.UserId)
}
