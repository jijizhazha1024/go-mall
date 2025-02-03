package logic

import (
	"context"
	"database/sql"
	"errors"

	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/internal/users_biz"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息方法
func (l *GetUserLogic) GetUser(in *users.GetUserRequest) (*users.GetUserResponse, error) {
	// todo: add your logic here and delete this line

	usermodel := user.NewUsersModel(l.svcCtx.Mysql)

	user, err := usermodel.FindOne(l.ctx, int64(in.UserId))
	if err != nil {
		if err == sql.ErrNoRows {
			return users_biz.HandleGetUsererror("user not found", 1, errors.New("user not found"))
		}
		return users_biz.HandleGetUsererror("sql error", 1, errors.New("sql error"))
	}
	if user.UserDeleted {
		return users_biz.HandleGetUsererror("user deleted", 1, errors.New("user deleted"))
	}

	return users_biz.HandleGetUserResp("get user success", 0, uint32(user.UserId), user.Username.String, user.Email.String)
}
