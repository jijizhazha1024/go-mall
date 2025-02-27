package logic

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

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
func GetCravatar(email string, size int, defaultImage, rating string, imgTag bool, attrs map[string]string) string {
	// 构建基本的 Cravatar URL
	baseURL := "https://cravatar.cn/avatar/"

	// 清理并对电子邮件地址进行 MD5 哈希处理
	email = strings.TrimSpace(strings.ToLower(email))
	hash := md5.New()
	hash.Write([]byte(email))
	emailHash := hex.EncodeToString(hash.Sum(nil))

	// 构建 Cravatar URL
	cravURL := fmt.Sprintf("%s%s?s=%d&d=%s&r=%s", baseURL, emailHash, size, defaultImage, rating)

	// 如果 imgTag 为 true，则返回完整的 <img> 标签
	if imgTag {
		imgTagStr := fmt.Sprintf(`<img src="%s"`, cravURL)
		for key, value := range attrs {
			imgTagStr += fmt.Sprintf(` %s="%s"`, key, value)
		}
		imgTagStr += " />"
		return imgTagStr
	}

	// 否则，仅返回 URL
	return cravURL
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

			// 获取头像
			size := 40
			defaultImage := "https://www.somewhere.com/homestar.jpg"
			rating := "g"
			imgTag := true
			attrs := map[string]string{"class": "avatar", "alt": "User Avatar"}
			avatar := GetCravatar(in.Email, size, defaultImage, rating, imgTag, attrs)

			// 加入布隆过滤器 在插入数据库之前防止数据库注册失败
			err = l.svcCtx.BF.Add([]byte(in.Email))
			if err != nil {
				l.Logger.Errorw("register bloom filter add failed", logx.Field("err", err),
					logx.Field("email", in.Email))
				return &users.RegisterResponse{}, err

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

				UserId:            uint32(userId),
				ActionType:        biz.Create,
				TargetTable:       "user",
				ActionDescription: "用户注册",
				TargetId:          int64(userId),
				ServiceName:       "users",
				ClientIp:          in.Ip,
			})
			if err != nil {
				l.Logger.Infow("register audit failed", logx.Field("err", err),
					logx.Field("email", in.Email))

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

			_, err = l.svcCtx.AuditRpc.CreateAuditLog(l.ctx, &audit.CreateAuditLogReq{

				UserId:            uint32(existUser.UserId),
				ActionType:        biz.Create,
				TargetTable:       "user",
				ActionDescription: "用户注册",
				TargetId:          int64(existUser.UserId),
				ServiceName:       "users",
				ClientIp:          in.Ip,
			})
			if err != nil {
				l.Logger.Infow("register audit failed", logx.Field("err", err),
					logx.Field("email", in.Email))

			}
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
