package logic

import (
	"context"
	"database/sql"
	"errors"

	"jijizhazha1024/go-mall/common/consts/code"
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

	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, int64(in.UserId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return &users.GetUserResponse{
				StatusCode: code.UserNotFound,
				StatusMsg:  code.UserNotFoundMsg,
			}, nil
		}
		logx.Errorw("find user failed", logx.Field("err", err), logx.Field("user_id", in.UserId))
		return &users.GetUserResponse{}, err

	}

	if user.UserDeleted {
		logx.Infow("you have deleted this user, please contact the administrator", logx.Field("user_id", in.UserId))
		return &users.GetUserResponse{
			StatusCode: code.UserInfoRetrievalFailed,
			StatusMsg:  code.UserInfoRetrievalFailedMsg,
		}, nil

	}
	return &users.GetUserResponse{

		UserId:    uint32(user.UserId),
		UserName:  user.Username.String,
		Email:     user.Email.String,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		LogoutAt:  user.LogoutAt.Time.Format("2006-01-02 15:04:05"),
		AvatarUrl: user.AvatarUrl.String,
	}, nil

}
