package logic

import (
	"context"
	"database/sql"
	"errors"

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
			l.Logger.Info("用户不存在：%d", in.UserId)
			return users_biz.HandleDeleteUsererror("用户不存在", 20016, err)
		}
		return users_biz.HandleDeleteUsererror("查询失败", 500, err)
	}
	// 删除用户
	if exituser.UserDeleted {
		l.Logger.Info("用户已删除", in.UserId)
		return users_biz.HandleDeleteUsererror("you have deleted this user", 20016, errors.New("you have deleted this user"))
	}
	err = l.svcCtx.UsersModel.UpdateDeletebyId(l.ctx, int64(in.UserId), true)
	if err != nil {
		return users_biz.HandleDeleteUsererror("删除失败", 20011, err)
	}

	return users_biz.HandleDeleteUserResp("删除成功", 0, in.UserId)
}
