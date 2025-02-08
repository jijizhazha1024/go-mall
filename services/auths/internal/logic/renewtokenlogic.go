package logic

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
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
	// parse jwt
	claims, err := token.ParseJWT(in.RefreshToken)
	var res = new(auths.AuthRenewalRes)
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
		logx.Infow("token expired", logx.Field("user_id", claims.UserID))
		return res, nil
	}
	// generate new jwt
	res.AccessToken, err = token.GenerateJWT(claims.UserID, claims.UserName, biz.TokenExpire)
	if err != nil {
		l.Logger.Errorw("access token generate failed",
			logx.Field("err", err),
			logx.Field("user_id", claims.UserID))
		return nil, err
	}
	res.RefreshToken, err = token.GenerateJWT(claims.UserID, claims.UserName, biz.TokenRenewalExpire)
	if err != nil {
		l.Logger.Errorw("refresh token generate failed",
			logx.Field("err", err),
			logx.Field("user_id", claims.UserID))
		return nil, err
	}
	return res, nil
}
