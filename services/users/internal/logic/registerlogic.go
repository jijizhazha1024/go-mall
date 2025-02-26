package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

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

func (l *RegisterLogic) GetAvatar() (string, error) {
	// 构建请求
	req, err := http.NewRequest("GET", "https://v2.xxapi.cn/api/head", nil)
	if err != nil {
		l.Logger.Infow("创建请求失败", logx.Field("err", err))
		return "", nil
	}

	// 设置 Header
	req.Header.Set("User-Agent", "xiaoxiaoapi/1.0.0 (https://xxapi.cn)")
	// 可选：其他 Header
	req.Header.Add("Accept", "application/json")
	// 发送请求
	resp, err := l.svcCtx.HttpClient.Do(req)
	if err != nil {
		l.Logger.Errorw("请求失败", logx.Field("err", err))
		return "", nil
	}
	defer resp.Body.Close()

	// 校验状态码
	if resp.StatusCode != http.StatusOK {
		l.Logger.Infow("非200响应", logx.Field("status", resp.Status))
		return "", nil
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Logger.Infow("读取响应失败", logx.Field("err", err))
		return "", nil
	}
	type ResponseURL struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}
	var responseURL ResponseURL
	err = json.Unmarshal(body, &responseURL)
	if err != nil {
		l.Logger.Infow("解析响应失败", logx.Field("err", err))
		return "", nil
	}

	return responseURL.Data, nil
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
			avatar, _ := l.GetAvatar()
			if avatar == "" {
				avatar = "https://www.gravatar.com/avatar/0000000000?d=mp"
			} //获取失败使用默认头像

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
