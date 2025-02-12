package logic

import (
	"context"

	"jijizhazha1024/go-mall/apis/user/internal/svc"
	"jijizhazha1024/go-mall/apis/user/internal/types"
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

	request := &usersclient.RegisterRequest{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	response, err := l.svcCtx.UserRpc.Register(l.ctx, request)
	if err != nil {
		return nil, err
	}
	authrespone, err := l.svcCtx.AuthsRpc.GenerateToken(l.ctx, &authsclient.AuthGenReq{
		UserId:   response.UserId,
		Username: "",
		ClientIp: "",
	})
	if err != nil {
		return nil, err
	}

	resp = &types.RegisterResponse{
		AccessToken:  authrespone.AccessToken,
		RefreshToken: authrespone.RefreshToken,
	}

	return resp, err
}
