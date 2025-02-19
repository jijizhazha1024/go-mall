package logic

import (
	"context"

	"github.com/zeromicro/x/errors"

	"jijizhazha1024/go-mall/apis/user/internal/svc"
	"jijizhazha1024/go-mall/apis/user/internal/types"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/auths/authsclient"
	"jijizhazha1024/go-mall/services/users/usersclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// todo: add your logic here and delete this line

	if req.Email == "" || req.Password == "" {
		return nil, errors.New(code.LoginMessageEmpty, code.LoginMessageEmptyMsg)
	}

	if req.Password != req.ConfirmPassword {
		l.Logger.Infow("密码不一致")
		return nil, errors.New(code.PasswordNotMatch, code.PasswordNotMatchMsg)

	}

	response, err := l.svcCtx.UserRpc.Register(l.ctx, &usersclient.RegisterRequest{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {

		l.Logger.Errorf("call rpc register failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, err.Error())
	} else {
		if response.StatusMsg != "" {
			l.Logger.Errorf("login failed", logx.Field("status_code", response.StatusCode), logx.Field("status_msg", response.StatusMsg))
			return nil, errors.New(int(response.StatusCode), response.StatusMsg)
		}

	}

	client_IP := l.ctx.Value(biz.ClientIPKey).(string)

	authrespone, err := l.svcCtx.AuthsRpc.GenerateToken(l.ctx, &authsclient.AuthGenReq{
		UserId:   response.UserId,
		Username: "",
		ClientIp: client_IP,
	})
	if err != nil {
		l.Logger.Errorf("call rpc generate token failed", logx.Field("err", err))
		return nil, errors.New(code.ServerError, code.ServerErrorMsg)

	}

	resp = &types.RegisterResponse{
		AccessToken:  authrespone.AccessToken,
		RefreshToken: authrespone.RefreshToken,
	}

	return resp, nil
}
