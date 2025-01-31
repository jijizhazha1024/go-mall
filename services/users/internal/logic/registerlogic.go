package logic

import (
	"context"
	"database/sql"
	"errors"

	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/users/internal/svc"
	"jijizhazha1024/go-mall/services/users/internal/users_biz"
	"jijizhazha1024/go-mall/services/users/users"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
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

// 注册方法
func (l *RegisterLogic) Register(in *users.RegisterRequest) (*users.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	//判断密码是否一致
	if in.Password != in.ConfirmPassword {
		l.Logger.Error("密码不一致")
		return users_biz.HandleRegisterResp("密码不一致", 0, 0, "token")

	}

	userMoel := user.NewUsersModel(l.svcCtx.Mysql)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Logger.Error("密码哈希生成失败", err)
		return users_biz.HandleRegisterResp("密码哈希生成失败", 0, 0, "token")

	}
	email := sql.NullString{
		String: in.Email,
		Valid:  true,
	}
	passwordhash := sql.NullString{
		String: string(hashedPassword),
		Valid:  true,
	}
	//判断邮箱是否已注册，如果已注册，是否处于删除状态
	existUser, err := userMoel.FindOneByEmail(l.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			l.Logger.Info("用户不存在", email)
			// 用户不存在，直接注册
			result, insertErr := userMoel.Insert(l.ctx, &user.Users{
				Email:        email,
				PasswordHash: passwordhash,
			})
			if insertErr != nil {
				l.Logger.Error("用户注册失败", insertErr)
				return users_biz.HandleRegistererror("用户注册失败", 1, errors.New("用户注册失败"))
			}
			userId, lastInsertErr := result.LastInsertId()
			if lastInsertErr != nil {
				l.Logger.Error("获取用户ID失败", lastInsertErr)
				return users_biz.HandleRegistererror("获取用户id失败", 1, errors.New("获取用户id失败"))
			}
			return users_biz.HandleRegisterResp("注册成功", 0, uint32(userId), "token")
		}
		l.Logger.Error("查询用户失败", err)
		return users_biz.HandleRegistererror("查询用户id失败", 1, errors.New("查询用户id失败"))
	}

	if existUser != nil {
		l.Logger.Info("用户已存在", existUser)
		// 用户已存在，判断是否处于删除状态
		userDeleted := existUser.UserDeleted
		if userDeleted { // 已删除
			// 将删除状态改为false
			updateErr := userMoel.UpdateDeletebyEmail(l.ctx, in.Email, false)
			if updateErr != nil {
				l.Logger.Error("更新用户状态失败", updateErr)
				return users_biz.HandleRegistererror("更新用户id失败", 1, errors.New("更新用户id失败"))
			}
			return users_biz.HandleRegisterResp("用户已存在，已恢复", 0, uint32(existUser.UserId), "token")
		} else { // 未删除
			l.Logger.Error("邮箱已注册")
			return users_biz.HandleRegistererror("邮箱已注册", 1, errors.New("邮箱已注册"))
		}

	}

	return users_biz.HandleRegisterResp("未知错误", 1, 0, "token")
}
