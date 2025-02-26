package logic

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"math/big"

	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/utils/cryptx"
	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/audit/audit"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var avatarList = []string{
	"http://example.com/avatar1.jpg",
	"http://example.com/avatar2.jpg",
	"http://example.com/avatar3.jpg",
	// 添加更多的头像URL
}

func getRandomAvatar() (string, error) {
	max := big.NewInt(int64(len(avatarList)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return avatarList[n.Int64()], nil
}

// 注册方法
func (l *RegisterLogic) Register(in *users.RegisterRequest) (*users.RegisterResponse, error) {
	// todo: add your logic here and delete this line

	email := sql.NullString{
		String: in.Email,
		Valid:  true,
	}
	PasswordHash := cryptx.PasswordEncrypt(in.Password)
	//判断邮箱是否已注册，如果已注册，是否处于删除状态
	existUser, err := l.svcCtx.UsersModel.FindOneByEmail(l.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			avatar, err := getRandomAvatar()
			if err != nil {
				l.Logger.Infow("register get avatar failed", logx.Field("err", err))
			}

			// 用户不存在，直接注册
			result, insertErr := l.svcCtx.UsersModel.Insert(l.ctx, &user.Users{
				Email:        email,
				PasswordHash: sql.NullString{String: PasswordHash, Valid: true},
				AvatarUrl:    sql.NullString{String: avatar, Valid: true},
			})

			if insertErr != nil {

				logx.Errorw("register insert user failed", logx.Field("err", insertErr), logx.Field("user_email", in.Email))
				return &users.RegisterResponse{}, err

			}

			userId, lastInsertErr := result.LastInsertId()
			if lastInsertErr != nil {
				l.Logger.Infow("register get user_id failed", logx.Field("err", lastInsertErr),
					logx.Field("email", in.Email))
				return &users.RegisterResponse{
					StatusCode: code.UserInfoRetrievalFailed,
					StatusMsg:  code.UserInfoRetrievalFailedMsg,
				}, nil

			}

			//审计操作
			_, err = l.svcCtx.AuditRpc.CreateAuditLog(l.ctx, &audit.CreateAuditLogReq{

				UserId:            uint32(existUser.UserId),
				ActionType:        biz.Create,
				TargetTable:       "user",
				ActionDescription: "用户注册",
				TargetId:          int64(userId),
				OldData:           "",
				NewData:           "",
				CreateAt:          0,
				ServiceName:       "users",
				ClientIp:          "",
			})
			if err != nil {
				l.Logger.Infow("register audit failed", logx.Field("err", err),
					logx.Field("email", in.Email))

				return &users.RegisterResponse{
					StatusCode: code.AuditRegisterFailed,
					StatusMsg:  code.AuditRegisterFailedMsg,
				}, nil
			}
			//埋点
			svc.UserRegCounter.Inc("success")

			return &users.RegisterResponse{

				UserId: uint32(userId),
			}, nil

		}

	}

	if existUser != nil {

		// 用户已存在，判断是否处于删除状态
		userDeleted := existUser.UserDeleted
		if userDeleted { // 已删除
			// 将删除状态改为false
			updateErr := l.svcCtx.UsersModel.UpdateDeletebyEmail(l.ctx, in.Email, false)
			if updateErr != nil {
				l.Logger.Errorw("register update user_deleted failed", logx.Field("err", updateErr),
					logx.Field("email", in.Email))
				return &users.RegisterResponse{}, updateErr

			}
			//给删除状态的用户 更新密码

			updatepasswordErr := l.svcCtx.UsersModel.UpdatePasswordHash(l.ctx, existUser.UserId, PasswordHash)
			if updatepasswordErr != nil {
				l.Logger.Errorw("register update password_hash failed", logx.Field("err", updatepasswordErr),
					logx.Field("email", in.Email))

				return &users.RegisterResponse{}, err

			}

			//审计操作
			// _, err = l.svcCtx.AuditRpc.CreateAuditLog(l.ctx, &audit.CreateAuditLogReq{

			// 	UserId:            uint32(existUser.UserId),
			// 	ActionType:        biz.Create,
			// 	TargetTable:       "user",
			// 	ActionDescription: "用户注册",
			// 	ServiceName:       "users",
			// })
			// if err != nil {
			// 	l.Logger.Infow("register audit failed", logx.Field("err", err),
			// 		logx.Field("email", in.Email))

			// 	return &users.RegisterResponse{
			// 		StatusCode: code.AuditRegisterFailed,
			// 		StatusMsg:  code.AuditRegisterFailedMsg,
			// 	}, nil
			// }
			//埋点操作
			svc.UserRegCounter.Inc("success")
			return &users.RegisterResponse{
				UserId: uint32(existUser.UserId),
			}, nil
		} else { // 未删除
			l.Logger.Infow("register  user already exist",
				logx.Field("email", in.Email))

			return &users.RegisterResponse{
				StatusCode: code.UserAlreadyExists,
				StatusMsg:  code.UserAlreadyExistsMsg,
			}, nil

		}

	}

	return &users.RegisterResponse{}, errors.New("register failed")

}
