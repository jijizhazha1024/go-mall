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

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户方法
func (l *UpdateUserLogic) UpdateUser(in *users.UpdateUserRequest) (*users.UpdateUserResponse, error) {
	// todo: add your logic here and delete this line

	update_user, err := l.svcCtx.UsersModel.FindOne(l.ctx, int64(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users_biz.HandleUpdateUsererror("user not found", 20016, errors.New("user not found"))
		}
		return users_biz.HandleUpdateUsererror("sql error", 500, errors.New("sql error"))
	}

	if update_user.UserDeleted {
		return users_biz.HandleUpdateUsererror("user deleted", 20016, errors.New("user deleted"))
	}

	email := sql.NullString{
		String: in.Email,
		Valid:  true,
	}

	var passworhash []byte
	if in.Password != "" { // 修改1: 处理密码为空字符串的情况
		passworhash, err = bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return users_biz.HandleUpdateUsererror("hash error", 1, errors.New("hash error"))
		}
	} else {
		// 如果密码为空，则不更新密码
		passworhash = nil
	}

	err = l.svcCtx.UsersModel.Update(l.ctx, &user.Users{
		UserId: int64(in.UserId),
		Username: sql.NullString{
			String: string(in.UsrName),
			Valid:  in.UsrName != "",
		},
		Email: email,
		PasswordHash: sql.NullString{
			String: string(passworhash),
			Valid:  passworhash != nil, // 修改1: 根据密码是否为空设置Valid字段
		},
	})
	if err != nil {
		return users_biz.HandleUpdateUsererror("sql error", 500, errors.New("sql error"))
	}
	return users_biz.HandleUpdateUserResp("user updated successfully", 0, in.UserId, "token") // 调用HandleUpdateUserResp方法返回响)

}
