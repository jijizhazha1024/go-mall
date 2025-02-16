package user_address

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserAddressesModel = (*customUserAddressesModel)(nil)

type (
	// UserAddressesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserAddressesModel.
	UserAddressesModel interface {
		userAddressesModel
		WithSession(session sqlx.Session) UserAddressesModel
		GetUserAddress(ctx context.Context, userId int32) (*UserAddresses, error)
		FindAllByUserId(ctx context.Context, userId int32) ([]*UserAddresses, error)
		DeleteByAddressIdandUserId(ctx context.Context, addressId int32, userId int32) error
		InsertWithSession(ctx context.Context, session sqlx.Session, data *UserAddresses) (sql.Result, error)
		GetUserAddressbyIdAndUserId(ctx context.Context, addressId int32, userId int32) (*UserAddresses, error)
		UpdateWithSession(ctx context.Context, session sqlx.Session, data *UserAddresses) (sql.Result, error)
		BatchUpdateDeFaultWithSession(ctx context.Context, session sqlx.Session, data []*UserAddresses) error
	}

	customUserAddressesModel struct {
		*defaultUserAddressesModel
	}
)

// NewUserAddressesModel returns a model for the database table.
func NewUserAddressesModel(conn sqlx.SqlConn) UserAddressesModel {
	return &customUserAddressesModel{
		defaultUserAddressesModel: newUserAddressesModel(conn),
	}
}

func (m *customUserAddressesModel) WithSession(session sqlx.Session) UserAddressesModel {
	return NewUserAddressesModel(sqlx.NewSqlConnFromSession(session))
}
func (m *customUserAddressesModel) GetUserAddress(ctx context.Context, userId int32) (*UserAddresses, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `is_default` = true limit 1", userAddressesRows, m.table)

	var resp UserAddresses
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserAddressesModel) FindAllByUserId(ctx context.Context, userId int32) ([]*UserAddresses, error) {
	query := fmt.Sprintf("select * from %s where `user_id` = ?", m.table)
	var resp []*UserAddresses

	err := m.conn.QueryRows(&resp, query, userId)
	return resp, err
}

func (m *customUserAddressesModel) DeleteByAddressIdandUserId(ctx context.Context, addressId int32, userId int32) error {
	query := fmt.Sprintf("delete from %s where `address_id` = ? and `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, addressId, userId)
	return err
}
func (m *customUserAddressesModel) GetUserAddressbyIdAndUserId(ctx context.Context, addressId int32, userId int32) (*UserAddresses, error) {
	query := fmt.Sprintf("select %s from %s where `address_id` = ? and `user_id` = ?", userAddressesRows, m.table)

	var resp UserAddresses
	err := m.conn.QueryRowCtx(ctx, &resp, query, addressId, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *customUserAddressesModel) BatchUpdateDeFaultWithSession(ctx context.Context, session sqlx.Session, data []*UserAddresses) error {
	for _, userAddress := range data {
		query := fmt.Sprintf("update %s set `is_default` = false where `user_id` = ?", m.table)
		_, err := session.ExecCtx(ctx, query, userAddress.UserId)
		if err != nil {
			return err
		}
	}
	return nil
}
func (m *customUserAddressesModel) InsertWithSession(ctx context.Context, session sqlx.Session, data *UserAddresses) (sql.Result, error) {
	// 定义插入的 SQL 语句
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, userAddressesRowsExpectAutoSet)
	// 使用 session 执行插入操作
	result, err := session.ExecCtx(ctx, query, data.UserId, data.DetailedAddress, data.City, data.Province, data.IsDefault, data.RecipientName, data.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (m *customUserAddressesModel) UpdateWithSession(ctx context.Context, session sqlx.Session, data *UserAddresses) (sql.Result, error) {
	// 定义更新的 SQL 语句
	query := fmt.Sprintf("update %s set %s where `address_id` = ?", m.table, userAddressesRowsWithPlaceHolder)
	// 使用 session 执行更新操作
	result, err := session.ExecCtx(ctx, query, data.UserId, data.DetailedAddress, data.City, data.Province, data.IsDefault, data.RecipientName, data.PhoneNumber, data.AddressId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
