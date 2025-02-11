package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/utils/token"
	"jijizhazha1024/go-mall/services/auths/auths"
	"jijizhazha1024/go-mall/services/auths/internal/svc"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuthenticationLogic) Authentication(in *auths.AuthReq) (*auths.AuthsRes, error) {
	res := new(auths.AuthsRes)
	// parse token
	claims, err := token.ParseJWT(in.Token)
	if err != nil {
		res.StatusCode = code.TokenValid
		res.StatusMsg = code.TokenInvalidMsg
		if errors.Is(err, jwt.ErrTokenExpired) {
			res.StatusCode = code.AuthExpired
			res.StatusMsg = code.AuthExpiredMsg
		}
		l.Logger.Infow("token parse failed",
			logx.Field("err", err),
			logx.Field("access_token", in.Token))
		return res, nil
	}
	clientIP := in.GetClientIp()
	if clientIP == "" {
		res.StatusCode = code.NotWithClientIP
		res.StatusMsg = code.NotWithClientIPMsg
		l.Logger.Infow("client ip is empty", logx.Field("access_token", in.Token))
		return res, nil
	}
	// check if the client IP has changed
	if clientIP != claims.ClientIP {
		res.StatusCode = code.AuthExpired
		res.StatusMsg = code.AuthExpiredMsg
		l.Logger.Infow("client ip changed",
			logx.Field("old_ip", claims.ClientIP),
			logx.Field("new_ip", clientIP))
		return res, nil
	}

	// comparison of jwt create time and user logout time
	logoutTime, err := l.svcCtx.UserModel.GetLogoutTime(l.ctx, int64(claims.UserID))
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		logx.Errorw("get logout time failed", logx.Field("err", err))
		return nil, err
	}
	issuedAt := claims.RegisteredClaims.IssuedAt
	if issuedAt.Before(logoutTime) {
		res.StatusCode = code.AuthExpiredByLogout
		res.StatusMsg = code.AuthExpiredByLogoutMsg
		// token expired
		logx.Infow("token expired by logout or re-login", logx.Field("user_id", claims.UserID),
			logx.Field("issued_at", issuedAt.Format("2006-01-02 15:04:05")), logx.Field("logout_time", logoutTime.Format("2006-01-02 15:04:05")))
		return res, nil
	}
	res.UserId = claims.UserID
	return res, nil
}
