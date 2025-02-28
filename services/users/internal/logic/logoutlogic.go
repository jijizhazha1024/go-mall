package logic

import (
	"context"
	"errors"
	"time"

	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/users/internal/svc"
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

			// 用户不存在
			return &users.LogoutResponse{
				StatusCode: code.UserNotFound,
				StatusMsg:  code.UserNotFoundMsg,
			}, nil

		}
		// 处理错误
		logx.Errorw("update logout time failed  query failed", logx.Field("err", err), logx.Field("user_id", in.UserId))
		return &users.LogoutResponse{}, err

	}
	logtoutime, err := l.svcCtx.UsersModel.GetLogoutTime(l.ctx, int64(in.UserId))
	if err != nil {

		// 处理错误
		logx.Errorw("get logout time failed  query failed", logx.Field("err", err), logx.Field("user_id", in.UserId))
		return &users.LogoutResponse{}, err
	}
	//审计操作

	return &users.LogoutResponse{

		LogoutTime: logtoutime.Unix(),
	}, nil

}
