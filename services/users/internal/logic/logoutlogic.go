package logic

import (
	"context"
	"errors"
	"time"

	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/internal/users_biz"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 登出方法
func (l *LogoutLogic) Logout(in *users.LogoutRequest) (*users.LogoutResponse, error) {

	// 在数据库中加入登出时间（这部分假设已经完成）
	logoutTime := time.Now()

	err := l.svcCtx.UsersModel.UpdateLogoutTime(l.ctx, int64(in.UserId), logoutTime)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			logx.Infof(code.UserNotFoundMsg, logx.Field("err", err),
				logx.Field("user id", in.UserId))

			// 用户不存在
			return users_biz.HandleLogoutUsererror(code.UserNotFoundMsg, code.UserNotFound, nil)
		}
		// 处理错误
		logx.Errorw(code.ServerErrorMsg, logx.Field("err", err), logx.Field("user id", in.UserId))
		return users_biz.HandleLogoutUsererror(code.ServerErrorMsg, code.ServerError, err)
	}

	// 从数据库中获取登出时间
	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, int64(in.UserId))
	if err != nil {
		logx.Infow(code.ServerErrorMsg)
		// 处理错误
		return users_biz.HandleLogoutUsererror(code.ServerErrorMsg, code.ServerError, err)
	}

	// 构造返回值

	return users_biz.HandleLogoutUserResp(code.LogoutSuccessMsg, code.LogoutSuccess, "token", user.LogoutAt.Time)
}
