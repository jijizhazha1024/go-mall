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
			logx.Infow("get user failed, user not found", logx.Field("err", err),
				logx.Field("user id", in.UserId))

			return users_biz.HandleGetUsererror(code.UserNotFoundMsg, code.UserNotFound, nil)
		}
		logx.Errorw(code.ServerErrorMsg, logx.Field("err", err), logx.Field("user_id", in.UserId))
		return users_biz.HandleGetUsererror(code.ServerErrorMsg, code.ServerError, err)
	}

	if user.UserDeleted {
		logx.Infow("you have deleted this user, please contact the administrator", logx.Field("user_id", in.UserId))

		return users_biz.HandleGetUsererror(code.UserInfoRetrievalFailedMsg, code.UserDeleted, nil)
	}

	return users_biz.HandleGetUserResp(code.UserInfoRetrievedMsg, code.UserInfoRetrieved, uint32(user.UserId), user.Username.String, user.Email.String, user.CreatedAt.Format("2006-01-02 15:04:05"), user.UpdatedAt.Format("2006-01-02 15:04:05"), user.LogoutAt.Time.Format("2006-01-02 15:04:05"), user.AvatarUrl.String)
}
