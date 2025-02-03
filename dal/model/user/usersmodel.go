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
		UpdateLogoutTime(ctx context.Context, userId int64, logoutTime time.Time) error
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

func (m *customUsersModel) UpdateDeletebyEmail(ctx context.Context, email string, userDeleted bool) error {
	query := fmt.Sprintf("update %s set `user_deleted` = ? where `email` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, userDeleted, email)
	return err
}
func (m *customUsersModel) FindAllEmails() ([]string, error) {
	query := fmt.Sprintf("SELECT email FROM %s", m.table)
	var emails []string
	err := m.conn.QueryRows(&emails, query)
	return emails, err
}
func (m *customUsersModel) UpdateLogoutTime(ctx context.Context, userId int64, logoutTime time.Time) error {
	query := fmt.Sprintf("update %s set `logout_at` = ? where `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, logoutTime, userId)
	return err
}

// 从数据库中获取登出时间
