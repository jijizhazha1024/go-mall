package logic

import (
	"context"
	"database/sql"
	"errors"

	"jijizhazha1024/go-mall/dal/model/user"
	"jijizhazha1024/go-mall/services/users/internal/svc"
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
		return nil, errors.New("密码不一致")
	}

	userMoel := user.NewUsersModel(l.svcCtx.Mysql)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Logger.Error("密码哈希生成失败", err)
		return nil, err
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
				return nil, errors.New("用户注册失败: " + insertErr.Error())
			}
			userId, lastInsertErr := result.LastInsertId()
			if lastInsertErr != nil {
				l.Logger.Error("获取用户ID失败", lastInsertErr)
				return nil, errors.New("获取用户ID失败")
			}
			return &users.RegisterResponse{
				StatusCode: 0,
				StatusMsg:  "注册成功",
				UserId:     uint32(userId),
				Token:      "token",
			}, nil
		}
		l.Logger.Error("查询用户失败", err)
		return nil, errors.New("查询用户失败: " + err.Error())
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
				return nil, errors.New("更新用户状态失败: " + updateErr.Error())
			}
			return &users.RegisterResponse{
				StatusCode: 0,
				StatusMsg:  "注册成功",
				UserId:     uint32(existUser.UserId),
				Token:      "token",
			}, nil
		} else { // 未删除
			l.Logger.Error("邮箱已注册")
			return nil, errors.New("邮箱已注册")
		}
	}

	return nil, errors.New("未知错误")
}
