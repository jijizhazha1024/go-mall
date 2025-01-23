package logic

import (
	"context"

	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除用户方法
func (l *DeleteUserLogic) DeleteUser(in *users.DeleteUserRequest) (*users.DeleteUserResponse, error) {
	// todo: add your logic here and delete this line

	return &users.DeleteUserResponse{}, nil
}
