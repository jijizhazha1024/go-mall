package logic

import (
	"context"
	"errors"
	"time"

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

	userMoel := user.NewUsersModel(l.svcCtx.Mysql)

	// 假设你有一个 userId 来标识用户

	// 在数据库中加入登出时间（这部分假设已经完成）
	logoutTime := time.Now()

	err := userMoel.UpdateLogoutTime(l.ctx, int64(in.UserId), logoutTime)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			// 用户不存在
			return users_biz.HandleLogoutUserResp("user not found", 1, 0, "", time.Now())
		}
		// 处理错误
		return users_biz.HandleLogoutUserResp("sql error", 1, 0, "", time.Now())
	}

	// 从数据库中获取登出时间
	user, err := userMoel.FindOne(l.ctx, int64(in.UserId))
	if err != nil {
		// 处理错误
		return users_biz.HandleLogoutUserResp("sql error", 1, 0, "", time.Now())
	}

	// 构造返回值

	return users_biz.HandleLogoutUserResp("success", 0, uint32(user.UserId), "token", logoutTime)
}
