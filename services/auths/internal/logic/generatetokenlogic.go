package logic

import (
	"context"
	"errors"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/utils/metadatactx"
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

	// 优先级：in.GetClientIp() > metadata 中的 biz.ClientIPKey > 返回错误
	clientIP := in.GetClientIp()
	if clientIP == "" {
		var ok bool
		clientIP, ok = metadatactx.ExtractFromMetadataCtx(l.ctx, biz.ClientIPKey)
		if !ok {
			l.Logger.Errorw("client ip is empty")
			return nil, errors.New("client ip is empty")
		}
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
	return &auths.AuthGenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
