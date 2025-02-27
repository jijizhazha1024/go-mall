package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/audit/audit"
	"jijizhazha1024/go-mall/services/users/internal/svc"
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
			return &users.UpdateUserResponse{
				StatusCode: code.UserNotFound,
				StatusMsg:  code.UserNotFoundMsg,
			}, nil

		}
		logx.Errorw(code.ServerErrorMsg, logx.Field("err", err), logx.Field("user_id", in.UserId))
		return &users.UpdateUserResponse{}, nil

	}

	if update_user.UserDeleted {

		logx.Infow(" update user have deleted", logx.Field("user_id", in.UserId), logx.Field("user_id", in.UserId))

		return &users.UpdateUserResponse{
			StatusCode: code.UserHaveDeleted,
			StatusMsg:  code.UserHaveDeletedMsg,
		}, nil
	}
	var username string
	username = in.UsrName
	var avatar_url string
	avatar_url = in.AvatarUrl
	if in.UsrName == "" {
		username = update_user.Username.String
	}
	if in.AvatarUrl == "" {
		avatar_url = update_user.AvatarUrl.String
	}
	err = l.svcCtx.UsersModel.UpdateUserNameandUrl(l.ctx, int64(in.UserId), username, avatar_url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logx.Infow("upate user not found", logx.Field("err", err),
				logx.Field("user_id", in.UserId))

			return &users.UpdateUserResponse{
				StatusCode: code.UserNotFound,
				StatusMsg:  code.UserNotFoundMsg,
			}, nil

		}

		return &users.UpdateUserResponse{}, err

	}

	//审计操作
	// 审计操作
	newData := ""
	if in.UsrName != "" && in.AvatarUrl != "" {
		newData = fmt.Sprintf("用户名: %s, 头像URL: %s", in.UsrName, in.AvatarUrl)
	} else if in.UsrName != "" {
		newData = fmt.Sprintf("用户名: %s", in.UsrName)
	} else if in.AvatarUrl != "" {
		newData = fmt.Sprintf("头像URL: %s", in.AvatarUrl)
	}

	_, err = l.svcCtx.AuditRpc.CreateAuditLog(l.ctx, &audit.CreateAuditLogReq{
		UserId:            uint32(in.UserId),
		ActionType:        biz.Update,
		TargetTable:       "user",
		ActionDescription: "用户更新",
		ServiceName:       "users",
		TargetId:          int64(in.UserId),
		OldData:           update_user.Username.String + ", " + update_user.AvatarUrl.String,
		NewData:           newData,
		ClientIp:          in.Ip,
	})
	if err != nil {
		logx.Infow("create audit log failed", logx.Field("err", err))

	}

	return &users.UpdateUserResponse{

		UserId:    in.UserId,
		AvatarUrl: in.AvatarUrl,

		UserName: in.UsrName,
	}, nil

}
