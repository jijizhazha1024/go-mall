package logic

import (
	"context"

	"jijizhazha1024/go-mall/apis/user/internal/svc"
	"jijizhazha1024/go-mall/apis/user/internal/types"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/users/usersclient"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteRequest) (resp *types.DeleteResponse, err error) {

	user_id := l.ctx.Value(biz.UserIDKey).(uint32)

	deleteresp, err := l.svcCtx.UserRpc.DeleteUser(l.ctx, &usersclient.DeleteUserRequest{
		UserId: uint32(user_id),
	})
	if err != nil {

		l.Logger.Errorf("call rpc deleteuser failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	} else {
		if deleteresp.StatusCode != code.UserDeleted {
			l.Logger.Errorf("delete failed", logx.Field("status_code", deleteresp.StatusCode), logx.Field("status_msg", deleteresp.StatusMsg))
			return nil, errors.New(int(deleteresp.StatusCode), deleteresp.StatusMsg)
		}

	}
	resp = &types.DeleteResponse{}

	return
}
