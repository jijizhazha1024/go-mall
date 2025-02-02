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
		return users_biz.HandleLoginerror("email or password is empty", 1, errors.New("email or password is empty"))
	}
	email := sql.NullString{
		String: in.Email,
		Valid:  true,
	}
	// 新增：布隆过滤器预检
	if !l.svcCtx.Bf.Contains(in.Email) {
		return users_biz.HandleLoginerror("email not allowed", 1, errors.New("email not allowed"))
	}
	// 2. 查询用户信息
	user, err := userMoel.FindOneByEmail(l.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users_biz.HandleLoginerror("user not found", 1, errors.New("user not found"))
		}
		return users_biz.HandleLoginerror("sql error", 1, errors.New("sql error"))
	}
	if user.UserDeleted {
		return users_biz.HandleLoginerror("user deleted", 1, errors.New("user deleted"))
	}

	// 3. 校验密码

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(in.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return users_biz.HandleLoginerror("password error", 1, errors.New("password error"))
		}
		return nil, err
	}

	return users_biz.HandleLoginResp("login success", 0, uint32(user.UserId), "", user.Username.String)
}
