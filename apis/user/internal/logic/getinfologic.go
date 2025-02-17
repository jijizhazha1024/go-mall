package logic

import (
	"context"
	"fmt"

	"jijizhazha1024/go-mall/apis/user/internal/svc"
	"jijizhazha1024/go-mall/apis/user/internal/types"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
)

type GetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInfoLogic {
	return &GetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInfoLogic) GetInfo(req *types.GetInfoRequest) (resp *types.GetInfoResponse, err error) {

	user_id := l.ctx.Value(biz.UserIDKey).(uint32)

	getresp, err := l.svcCtx.UserRpc.GetUser(l.ctx, &users.GetUserRequest{
		UserId: user_id,
	})
	if err != nil {

		l.Logger.Errorf("call rpc getuser failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	} else {
		if getresp.StatusCode != code.UserInfoRetrieved {
			l.Logger.Errorf("login failed", logx.Field("status_code", getresp.StatusCode), logx.Field("status_msg", getresp.StatusMsg))
			return nil, errors.New(int(getresp.StatusCode), getresp.StatusMsg)
		}

	}
	resp = &types.GetInfoResponse{
		UserId:    int64(getresp.UserId),
		LogoutAt:  getresp.LogoutAt,
		CreatedAt: getresp.CreatedAt,
		UpdateAt:  getresp.UpdatedAt,
		Email:     getresp.Email,
		UserName:  getresp.UserName,
	}
	fmt.Println("resp:", resp)

	return resp, nil
}
