package logic

import (
	"context"
	"database/sql"
	"errors"

	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 登录方法
func (l *LoginLogic) Login(in *users.LoginRequest) (*users.LoginResponse, error) {
	// todo: add your logic here and delete this line

	userMoel := user.NewUsersModel(l.svcCtx.Mysql)
	// 1. 校验参数
	if in.Email == "" || in.Password == "" {
		return &users.LoginResponse{
			StatusCode: 1,
			StatusMsg:  "email or password is empty",
		}, nil
	}
	email := sql.NullString{
		String: in.Email,
		Valid:  true,
	}
	// 2. 查询用户信息
	user, err := userMoel.FindOneByEmail(l.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &users.LoginResponse{
				StatusCode: 1,
				StatusMsg:  "user not found",
			}, nil
		}
		return nil, err
	}

	// 3. 校验密码

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(in.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return &users.LoginResponse{
				StatusCode: 1,
				StatusMsg:  "password is incorrect",
			}, nil
		}
		return nil, err
	}

	return &users.LoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     uint32(user.UserId),
		Token:      "token",
		UserName:   user.Username.String,
	}, nil
}
