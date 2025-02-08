package logic

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/utils/token"
	"jijizhazha1024/go-mall/services/auths/auths"
	"jijizhazha1024/go-mall/services/auths/internal/svc"
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
		return res, nil
	}
	// comparison of jwt create time and user logout time
	logOutTime := int64(0)
	if claims.RegisteredClaims.IssuedAt.Unix() <= logOutTime {
		res.StatusCode = code.AuthExpired
		res.StatusMsg = code.AuthExpiredMsg
		// token expired
		logx.Infow("token expired by logout", logx.Field("user_id", claims.UserID))
		return res, nil
	}
	res.UserId = claims.UserID
	return res, nil
}
