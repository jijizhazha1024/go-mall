package logic

import (
	"context"
	"fmt"

	"github.com/zeromicro/x/errors"

	"jijizhazha1024/go-mall/apis/user/internal/svc"
	"jijizhazha1024/go-mall/apis/user/internal/types"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/auths/authsclient"
	"jijizhazha1024/go-mall/services/users/usersclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line

	if req.Email == "" || req.Password == "" {
		return nil, errors.New(code.LoginMessageEmpty, code.LoginMessageEmptyMsg)
	}

	loginres, err := l.svcCtx.UserRpc.Login(l.ctx, &usersclient.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {

		l.Logger.Errorw("call rpc login failed", logx.Field("err", err))
		fmt.Println("loginres:", loginres)
		fmt.Println("err:", err)
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)
	} else if loginres.StatusMsg != "" {

		return nil, errors.New(int(loginres.StatusCode), loginres.StatusMsg)

	}

	client_IP := l.ctx.Value(biz.ClientIPKey).(string)

	authrespone, err := l.svcCtx.AuthsRpc.GenerateToken(l.ctx, &authsclient.AuthGenReq{
		UserId:   loginres.UserId,
		Username: loginres.UserName,
		ClientIp: client_IP,
	})
	if err != nil {
		l.Logger.Errorw("call rpc  auth token failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)

	}

	resp = &types.LoginResponse{
		AccessToken:  authrespone.AccessToken,
		RefreshToken: authrespone.RefreshToken,
	}

	return resp, nil
}
