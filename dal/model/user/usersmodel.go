package user

import (
	"context"
	"fmt"

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
