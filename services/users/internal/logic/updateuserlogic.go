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
			logx.Infow("update user not found", logx.Field("err", err),
				logx.Field("user_id", in.UserId))

			return users_biz.HandleUpdateUsererror(code.UserNotFoundMsg, code.UserNotFound, nil)
		}
		logx.Errorw(code.ServerErrorMsg, logx.Field("err", err), logx.Field("user_id", in.UserId))

		return users_biz.HandleUpdateUsererror(code.ServerErrorMsg, code.ServerError, err)
	}

	if update_user.UserDeleted {

		logx.Infow(" update user have deleted", logx.Field("user_id", in.UserId), logx.Field("user_id", in.UserId))

		return users_biz.HandleUpdateUsererror(code.UserHaveDeletedMsg, code.UserNotFound, nil)
	}

	err = l.svcCtx.UsersModel.UpdateUserName(l.ctx, int64(in.UserId), in.UsrName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logx.Infow("upate user not found",
				logx.Field("user_id", in.UserId))

			return users_biz.HandleUpdateUsererror(code.UserNotFoundMsg, code.UserNotFound, nil)
		}
		logx.Errorw(code.ServerErrorMsg, logx.Field("err", err), logx.Field("user id", in.UserId))
		return users_biz.HandleUpdateUsererror(code.ServerErrorMsg, code.ServerError, err)
	}
	return users_biz.HandleUpdateUserResp(code.UserUpdatedMsg, code.UserUpdated, in.UserId, in.UsrName) // 调用HandleUpdateUserResp方法返回响)

}
