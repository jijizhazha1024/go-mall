package logic

import (
	"context"

	"jijizhazha1024/go-mall/apis/user/internal/svc"
	"jijizhazha1024/go-mall/apis/user/internal/types"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateRequest) (resp *types.UpdateResponse, err error) {
	// todo: add your logic here and delete this line

	updateresp, err := l.svcCtx.UserRpc.UpdateUser(l.ctx, &users.UpdateUserRequest{
		Password: req.Password,
		UserId:   uint32(req.UserId),
		Email:    req.Email,
		UsrName:  req.UserName,
	})

	if err != nil {

		l.Logger.Errorf("call rpc update failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	} else {
		if updateresp.StatusCode != code.UserCreated {
			l.Logger.Errorf("update failed", logx.Field("status_code", updateresp.StatusCode), logx.Field("status_msg", updateresp.StatusMsg))
			return nil, errors.New(int(updateresp.StatusCode), updateresp.StatusMsg)
		}

	}

	resp = &types.UpdateResponse{

		Email:    updateresp.Email,
		UserName: updateresp.UserName,
		UserId:   int64(updateresp.UserId),
	}

	return
}
