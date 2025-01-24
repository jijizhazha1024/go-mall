package logic

import (
	"context"
	"database/sql"

	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 注册方法
func (l *RegisterLogic) Register(in *users.RegisterRequest) (*users.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	userMoel := user.NewUsersModel(l.svcCtx.Mysql)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Logger.Error("密码哈希生成失败", err)
		return nil, err
	}
	email := sql.NullString{
		String: in.Email,
		Valid:  true,
	}
	passwordhash := sql.NullString{
		String: string(hashedPassword),
		Valid:  true,
	}
	_, err = userMoel.Insert(l.ctx, &user.Users{
		Email:        email,
		PasswordHash: passwordhash,
	})
	if err != nil {
		l.Logger.Error("用户注册失败", err)
		return nil, err
	}

	return &users.RegisterResponse{
		StatusCode: 200,
		StatusMsg:  "注册成功",
		Data: &users.RegisterResponseData{
			UserId: 1,

			Token: "token",
		},
	}, nil
}
