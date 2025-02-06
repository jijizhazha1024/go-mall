package logic

import (
	"context"
	"database/sql"

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
		if err == sql.ErrNoRows {
			logx.Error(code.UserNotFoundMsg, user.UserId, err)
			return users_biz.HandleGetUsererror(code.UserNotFoundMsg, code.UserNotFound)
		}
		logx.Error(code.ServerErrorMsg, err)
		return users_biz.HandleGetUsererror(code.ServerErrorMsg, code.ServerError)
	}

	if user.UserDeleted {
		logx.Error(code.UserInfoRetrievalFailedMsg, user.UserId, err)
		return users_biz.HandleGetUsererror(code.UserInfoRetrievalFailedMsg, code.UserDeleted)
	}

	return users_biz.HandleGetUserResp(code.UserInfoRetrievedMsg, code.UserInfoRetrieved, uint32(user.UserId), user.Username.String, user.Email.String)
}
