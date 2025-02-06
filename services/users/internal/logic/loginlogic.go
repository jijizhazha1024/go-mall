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

	// 1. 校验参数
	if in.Email == "" || in.Password == "" {
		return users_biz.HandleLoginerror("email or password is empty", 400, errors.New("email or password is empty"))
	}
	email := sql.NullString{
		String: in.Email,
		Valid:  true,
	}
	// 新增：布隆过滤器预检
	if !l.svcCtx.Bf.Contains(in.Email) {
		return users_biz.HandleLoginerror(code.UserNotFoundMsg, code.UserNotFound, nil)
	}
	// 2. 查询用户信息
	user, err := l.svcCtx.UsersModel.FindOneByEmail(l.ctx, email)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			logx.Error(code.UserNotFoundMsg)
			logx.Field("err", err)
			logx.Field("user id", in.Email)
			return users_biz.HandleLoginerror(code.UserNotFoundMsg, code.UserNotFound, nil)
		}
		logx.Error(code.ServerErrorMsg, err)
		return users_biz.HandleLoginerror(code.ServerErrorMsg, code.ServerError, err)
	}
	if user.UserDeleted {
		logx.Error(code.UserHaveDeletedMsg, user.Email, err)
		logx.Field("err", err)
		logx.Field("email", user.Email)
		return users_biz.HandleLoginerror(code.UserHaveDeletedMsg, code.UserHaveDeleted, nil)
	}

	// 3. 校验密码

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(in.Password))
	if err != nil {
		logx.Error(code.LoginFailedMsg, user.Email, err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return users_biz.HandleLoginerror("password error", 400, nil)

		}
		return nil, err
	}

	return users_biz.HandleLoginResp(code.LoginSuccessMsg, code.LoginSuccess, uint32(user.UserId), "", user.Username.String)
}
