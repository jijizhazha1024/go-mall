package logic

import (
	"context"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/utils/token"
	"jijizhazha1024/go-mall/services/auths/auths"
	"jijizhazha1024/go-mall/services/auths/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GenerateToken 生成toke
func (l *GenerateTokenLogic) GenerateToken(in *auths.AuthGenReq) (*auths.AuthGenRes, error) {
	res := new(auths.AuthGenRes)
	clientIP := in.GetClientIp()
	if clientIP == "" {
		res.StatusCode = code.NotWithClientIP
		res.StatusMsg = code.NotWithClientIPMsg
		l.Logger.Infow("client ip is empty", logx.Field("user_id", in.UserId))
		return res, nil
	}
	// 生成access token
	accessToken, err := token.GenerateJWT(in.UserId, in.Username, clientIP, biz.TokenExpire)
	if err != nil {
		l.Logger.Errorw("access token generate failed",
			logx.Field("err", err),
			logx.Field("client_ip", clientIP),
			logx.Field("user_id", in.UserId))
		return nil, err
	}
	// 生成refresh token
	refreshToken, err := token.GenerateJWT(in.UserId, in.Username, clientIP, biz.TokenRenewalExpire)
	if err != nil {
		l.Logger.Errorw("refresh token generate failed",
			logx.Field("err", err),
			logx.Field("client_ip", clientIP),
			logx.Field("user_id", in.UserId))
		return nil, err
	}
	// 返回access token和refresh token
	l.Logger.Infow("tokens generated successfully",
		logx.Field("user_id", in.UserId),
		logx.Field("client_ip", clientIP))
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	return res, nil
}
