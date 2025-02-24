package user_address

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserAddressesModel = (*customUserAddressesModel)(nil)

type (
	// UserAddressesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserAddressesModel.
	UserAddressesModel interface {
		userAddressesModel
		GetUserAddressExistsByIdAndUserId(ctx context.Context, addressId int32, userId int32) (bool, error)
		WithSession(session sqlx.Session) UserAddressesModel

		FindAllByUserId(ctx context.Context, userId int32) ([]*UserAddresses, error)
		DeleteByAddressIdandUserId(ctx context.Context, addressId int32, userId int32) error
		InsertWithSession(ctx context.Context, session sqlx.Session, data *UserAddresses) (sql.Result, error)
		GetUserAddressbyIdAndUserId(ctx context.Context, addressId int32, userId int32) (*UserAddresses, error)
		UpdateWithSession(ctx context.Context, session sqlx.Session, data *UserAddresses) (sql.Result, error)

		BatchUpdateDeFaultWithSession(ctx context.Context, session sqlx.Session, data []*UserAddresses) error
	}

	customUserAddressesModel struct {
		cacheConf cache.CacheConf
		*defaultUserAddressesModel
	}
)

// NewUserAddressesModel returns a model for the database table.
func NewUserAddressesModel(conn sqlx.SqlConn, c cache.CacheConf) UserAddressesModel {
	return &customUserAddressesModel{
		defaultUserAddressesModel: newUserAddressesModel(conn, c),
		cacheConf:                 c,
	}
}

// 修改WithSession方法（原代码存在未定义的c变量）
func (m *customUserAddressesModel) WithSession(session sqlx.Session) UserAddressesModel {
	return NewUserAddressesModel(
		sqlx.NewSqlConnFromSession(session),
		m.cacheConf, // 使用保存的缓存配置
	)
}

func (m *customUserAddressesModel) FindAllByUserId(ctx context.Context, userId int32) ([]*UserAddresses, error) {
	var resp []*UserAddresses

	// 直接使用原始数据库连接（不走缓存）
	err := m.defaultUserAddressesModel.CachedConn.QueryRowsNoCacheCtx(
		ctx,
		&resp,
		"SELECT "+userAddressesRows+" FROM "+m.table+" WHERE `user_id` = ?",
		userId,
	)

	switch {
	case err == nil:
		return resp, nil
	case err == sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserAddressesModel) DeleteByAddressIdandUserId(ctx context.Context, addressId int32, userId int32) error {
	// 添加全量缓存清除
	keys := []string{
		fmt.Sprintf("userAddress:%d:%d", userId, addressId), // 主键缓存
	}

	// 带缓存的删除操作
	_, err := m.CachedConn.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		query := fmt.Sprintf("delete from %s where `address_id` = ? and `user_id` = ?", m.table)
		return conn.ExecCtx(ctx, query, addressId, userId)
	}, keys...)

	return err
}
func (m *customUserAddressesModel) GetUserAddressExistsByIdAndUserId(ctx context.Context, addressId int32, userId int32) (bool, error) {
	cacheKey := fmt.Sprintf("userAddressExists:%d:%d", userId, addressId)
	var exists int8

	// 使用联合主键缓存
	err := m.CachedConn.QueryRowCtx(ctx, &exists, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select exists(select 1 from %s where `address_id` = ? and `user_id` = ?)",
			m.table)
		return conn.QueryRowCtx(ctx, v, query, addressId, userId)
	})

	switch err {
	case nil:
		return exists == 1, nil
	case sqlx.ErrNotFound:
		return false, ErrNotFound
	default:
		return false, err
	}
}
func (m *customUserAddressesModel) GetUserAddressbyIdAndUserId(ctx context.Context, addressId int32, userId int32) (*UserAddresses, error) {
	cacheKey := fmt.Sprintf("userAddress:%d:%d", userId, addressId)
	var resp UserAddresses

	// 使用联合主键缓存
	err := m.CachedConn.QueryRowCtx(ctx, &resp, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `address_id` = ? and `user_id` = ?",
			userAddressesRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, addressId, userId)
	})

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
	// 事务成功后，批量删除相关用户的默认地址缓存

	var keys []string
	for _, addr := range data {
		keys = append(keys, fmt.Sprintf("userAddress:%d:%d", addr.UserId, addr.AddressId))
	}
	err := m.CachedConn.DelCacheCtx(ctx, keys...)
	if err != nil {
		// 可选择回滚事务或记录日志
		return err
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
