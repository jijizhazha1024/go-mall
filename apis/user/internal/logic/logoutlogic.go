package logic

import (
	"context"

	"jijizhazha1024/go-mall/apis/user/internal/svc"
	"jijizhazha1024/go-mall/apis/user/internal/types"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/users/usersclient"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutRequest) (resp *types.LogoutResponse, err error) {
	// todo: add your logic here and delete this line
	logoutrep, err := l.svcCtx.UserRpc.Logout(l.ctx, &usersclient.LogoutRequest{

		UserId: uint32(req.UserId),
	})
	if err != nil {

		l.Logger.Errorf("call rpc logout failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	} else {
		if logoutrep.StatusCode != code.UserCreated {
			l.Logger.Errorf("logout failed", logx.Field("status_code", logoutrep.StatusCode), logx.Field("status_msg", logoutrep.StatusMsg))
			return nil, errors.New(int(logoutrep.StatusCode), logoutrep.StatusMsg)
		}
	}

	resp = &types.LogoutResponse{
		Logout_at: logoutrep.LogoutTime,
	}

	return
}
