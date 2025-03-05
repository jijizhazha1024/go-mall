package user

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		withSession(session sqlx.Session) UsersModel
		UpdateDeletebyId(ctx context.Context, userId int64, userDeleted bool) error
		UpdateDeletebyEmail(ctx context.Context, email string, userDeleted bool) error
		FindAllEmails() ([]string, error)
		UpdateUserNameandUrl(ctx context.Context, userId int64, userName string, AvatarUrl string) error
		GetLogoutTime(ctx context.Context, userId int64) (time.Time, error)
		UpdateLoginTime(ctx context.Context, userId int64, loginTime time.Time) error
		UpdateLogoutTime(ctx context.Context, userId int64, logoutTime time.Time) error
		GetLoginTime(ctx context.Context, userId int64) (time.Time, error)
		UpdatePasswordHash(ctx context.Context, userId int64, passwordHash string) error
		// 从数据库中获取登出时间

	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn),
	}
}

func (m *customUsersModel) withSession(session sqlx.Session) UsersModel {
	return NewUsersModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customUsersModel) UpdateDeletebyId(ctx context.Context, userId int64, userDeleted bool) error {
	query := fmt.Sprintf("update %s set `user_deleted` = ? where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userDeleted, userId)
	return err
}
func (m *customUsersModel) UpdateUserNameandUrl(ctx context.Context, userId int64, userName string, AvatarUrl string) error {
	query := fmt.Sprintf("update %s set `username` = ?, `avatar_url` = ? where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userName, AvatarUrl, userId)
	return err
}

func (m *customUsersModel) UpdateDeletebyEmail(ctx context.Context, email string, userDeleted bool) error {
	query := fmt.Sprintf("update %s set `user_deleted` = ? where `email` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userDeleted, email)
	return err
}

func (m *customUsersModel) UpdateLogoutTime(ctx context.Context, userId int64, logoutTime time.Time) error {
	query := fmt.Sprintf("update %s set `logout_at` = ? where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, logoutTime, userId)
	return err
}

func (m *customUsersModel) UpdateLoginTime(ctx context.Context, userId int64, loginTime time.Time) error {
	query := fmt.Sprintf("update %s set `login_at` = ? where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, loginTime, userId)
	return err
}

func (m *customUsersModel) FindAllEmails() ([]string, error) {
	query := fmt.Sprintf("SELECT email FROM %s", m.table)
	var emails []string
	err := m.conn.QueryRows(&emails, query)
	return emails, err
}

func (m *customUsersModel) GetLoginTime(ctx context.Context, userId int64) (time.Time, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", usersRows, m.table)
	var user Users

	err := m.conn.QueryRowCtx(ctx, &user, query, userId)
	t := time.Time{}
	switch err {
	case nil:
		return user.LoginAt.Time, nil
	default:
		return t, err
	}

}

func (m *customUsersModel) GetLogoutTime(ctx context.Context, userId int64) (time.Time, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? limit 1", usersRows, m.table)
	var user Users
	now := time.Now()
	err := m.conn.QueryRowCtx(ctx, &user, query, userId)
	switch err {
	case nil:
		return user.LogoutAt.Time, nil
	case sqlx.ErrNotFound:
		return now.Add(2 * time.Hour), ErrNotFound
	default:
		return time.Time{}, err
	}
}
func (m *customUsersModel) UpdatePasswordHash(ctx context.Context, userId int64, passwordHash string) error {
	query := fmt.Sprintf("update %s set `password_hash` = ? where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, passwordHash, userId)
	return err
}

// 从数据库中获取登出时间
