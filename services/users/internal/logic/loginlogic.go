package logic

import (
	"context"
	"database/sql"
	"errors"

	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/utils/cryptx"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
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

	email := sql.NullString{
		String: in.Email,
		Valid:  true,
	}

	//bf
	bloom_exist, err := l.svcCtx.BF.Exists([]byte(in.Email))
	if err != nil {
		logx.Errorw("login failed, bloom filter query failed",
			logx.Field("err", err),
			logx.Field("user email", in.Email),
		)
		return &users.LoginResponse{}, err
	}
	if !bloom_exist {
		logx.Infow("login failed, bloom filter not exist", logx.Field("email", in.Email))

		return &users.LoginResponse{
			StatusCode: code.UserNotFound,
			StatusMsg:  code.UserNotFoundMsg,
		}, nil
	}

	// 2. 查询用户信息
	user, err := l.svcCtx.UsersModel.FindOneByEmail(l.ctx, email)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			logx.Infow("login failed, user not found", logx.Field("err", err),
				logx.Field("email", in.Email))
			return &users.LoginResponse{
				StatusCode: code.UserNotFound,
				StatusMsg:  code.UserNotFoundMsg,
			}, nil

		}
		logx.Errorw("login failed, database query failed",
			logx.Field("err", err),
			logx.Field("user email", in.Email),
		)
		return &users.LoginResponse{}, err
	}
	if user.UserDeleted {
		logx.Infow("login failed, user have deleted", logx.Field("email", user.Email))

		return &users.LoginResponse{
			StatusCode: code.UserHaveDeleted,
			StatusMsg:  code.UserHaveDeletedMsg,
		}, nil
	}

	// 3. 校验密码
	if !cryptx.PasswordVerify(in.Password, user.PasswordHash.String) {
		logx.Infow("login failed, password not match")

		return &users.LoginResponse{
			StatusCode: code.PasswordNotMatch,
			StatusMsg:  code.PasswordNotMatchMsg,
		}, nil
	}

	//审计操作

	return &users.LoginResponse{

		UserId: uint32(user.UserId),

		UserName: user.Username.String,
	}, nil

}
