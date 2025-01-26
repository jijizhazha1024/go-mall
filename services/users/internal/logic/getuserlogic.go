package logic

import (
	"context"
	"database/sql"
	"errors"

	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/users/internal/svc"
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
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &users.GetUserResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     uint32(user.UserId),
		Email:      user.Email.String,
		UserName:   user.Username.String,
	}, nil
}
