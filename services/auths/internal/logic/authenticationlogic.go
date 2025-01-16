package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"strconv"

	"jijizhazha1024/go-mall/services/auths/auths"
	"jijizhazha1024/go-mall/services/auths/internal/svc"

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
	// valid
	key := fmt.Sprintf(biz.TokenPrefixKey, in.Token)
	val, err := l.svcCtx.Rdb.GetCtx(l.ctx, key)
	if err != nil {
		// key is not exist err
		if errors.Is(err, redis.Nil) {
			res.StatusCode = code.AuthExpired
			res.StatusMsg = code.AuthExpiredMsg
			l.Logger.Infow("token is valid or expired", logx.Field("token", in.Token))
			return res, nil
		}
		l.Logger.Errorw("redis get error", logx.Field("token", in.Token), logx.Field("err", err))
		return nil, err
	}
	if val == "" {
		res.StatusCode = code.AuthExpired
		res.StatusMsg = code.AuthExpiredMsg
		l.Logger.Infow("token is valid or expired", logx.Field("token", in.Token))
		return res, nil
	}
	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		l.Logger.Errorw("strconv.ParseInt error",
			logx.Field("token", in.Token),
			logx.Field("val", val),
			logx.Field("err", err),
		)
		return nil, err
	}
	res.UserId = uint32(userID)
	return res, nil
}
