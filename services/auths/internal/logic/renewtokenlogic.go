package logic

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/utils/metadatactx"
	"jijizhazha1024/go-mall/common/utils/token"

	"jijizhazha1024/go-mall/services/auths/auths"
	"jijizhazha1024/go-mall/services/auths/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RenewTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRenewTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RenewTokenLogic {
	return &RenewTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RenewToken 续期身份
func (l *RenewTokenLogic) RenewToken(in *auths.AuthRenewalReq) (*auths.AuthRenewalRes, error) {
	res := new(auths.AuthRenewalRes)

	// parse jwt
	claims, err := token.ParseJWT(in.RefreshToken)
	if err != nil {
		res.StatusCode = code.TokenValid
		res.StatusMsg = code.TokenInvalidMsg
		if errors.Is(err, jwt.ErrTokenExpired) {
			res.StatusCode = code.AuthExpired
			res.StatusMsg = code.AuthExpiredMsg
		}
		l.Logger.Infow("token parse failed",
			logx.Field("err", err),
			logx.Field("refresh_token", in.RefreshToken))
		return res, nil
	}
	// comparison of jwt create time and user logout time
	logoutTime, err := l.svcCtx.UserModel.GetLogoutTime(l.ctx, int64(claims.UserID))
	if err != nil {
		logx.Errorw("get logout time failed", logx.Field("err", err))
		return nil, err
	}
	issuedAt := claims.RegisteredClaims.IssuedAt
	if issuedAt.Before(logoutTime) {
		res.StatusCode = code.AuthExpired
		res.StatusMsg = code.AuthExpiredMsg
		// token expired
		logx.Infow("token expired by logout or re-login",
			logx.Field("user_id", claims.UserID),
			logx.Field("issued_at", issuedAt.Format("2006-01-02 15:04:05")),
			logx.Field("logout_time", logoutTime.Format("2006-01-02 15:04:05")))
		return res, nil
	}
	// 获取客户端IP
	clientIP := in.GetClientIp()
	if clientIP == "" {
		var ok bool
		clientIP, ok = metadatactx.ExtractFromMetadataCtx(l.ctx, biz.ClientIPKey)
		if !ok {
			l.Logger.Infow("client ip is empty", logx.Field("user_id", claims.UserID))
			return nil, errors.New("client ip is empty")
		}
	}
	// generate new jwt
	res.AccessToken, err = token.GenerateJWT(claims.UserID, claims.UserName, clientIP, biz.TokenExpire)
	if err != nil {
		l.Logger.Errorw("access token generate failed",
			logx.Field("err", err),
			logx.Field("user_id", claims.UserID))
		return nil, err
	}
	res.RefreshToken, err = token.GenerateJWT(claims.UserID, claims.UserName, clientIP, biz.TokenRenewalExpire)
	if err != nil {
		l.Logger.Errorw("refresh token generate failed",
			logx.Field("err", err),
			logx.Field("user_id", claims.UserID))
		return nil, err
	}
	l.Logger.Infow("tokens renewed successfully",
		logx.Field("user_id", claims.UserID),
		logx.Field("client_ip", clientIP))
	return res, nil
}
