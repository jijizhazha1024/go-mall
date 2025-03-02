package checkout

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CheckoutsModel = (*customCheckoutsModel)(nil)

type (
	// CheckoutsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCheckoutsModel.
	CheckoutsModel interface {
		checkoutsModel
		withSession(session sqlx.Session) CheckoutsModel
		UpdateStatusWithSession(ctx context.Context, session sqlx.Session, status int64, userId int32, preOrderId string) error
		FindOneByUserIdAndPreOrderIdWithSession(ctx context.Context, session sqlx.Session, userId int32, preOrderId string) (*Checkouts, error)
		CountByUserId(ctx context.Context, userId uint32) (int64, error)
		FindOneByUserIdAndPreOrderId(ctx context.Context, userId int32, preOrderId string) (*Checkouts, error)
		FindByUserId(ctx context.Context, userId uint32, page int32, pageSize int32) ([]*Checkouts, error)
	}

	customCheckoutsModel struct {
		*defaultCheckoutsModel
	}
)

// NewCheckoutsModel returns a model for the database table.
func NewCheckoutsModel(conn sqlx.SqlConn) CheckoutsModel {
	return &customCheckoutsModel{
		defaultCheckoutsModel: newCheckoutsModel(conn),
	}
}

func (m *customCheckoutsModel) withSession(session sqlx.Session) CheckoutsModel {
	return NewCheckoutsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customCheckoutsModel) UpdateStatusWithSession(ctx context.Context, session sqlx.Session, status int64, userId int32, preOrderId string) error {
	query := `UPDATE checkout_orders 
			  SET status = ? 
			  WHERE user_id = ? AND pre_order_id = ?`
	_, err := session.ExecCtx(ctx, query, status, userId, preOrderId)
	return err
}

func (m *customCheckoutsModel) FindOneByUserIdAndPreOrderIdWithSession(ctx context.Context, session sqlx.Session, userId int32, preOrderId string) (*Checkouts, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `user_id` = ? AND `pre_order_id` = ? LIMIT 1 FOR SHARE", checkoutsRows, m.table)

	var resp Checkouts
	err := session.QueryRowCtx(ctx, &resp, query, userId, preOrderId) // 使用 session 查询
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, sqlx.ErrNotFound
		}
		return nil, err
	}
	return &resp, nil
}
func (m *customCheckoutsModel) FindOneByUserIdAndPreOrderId(ctx context.Context, userId int32, preOrderId string) (*Checkouts, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `pre_order_id` = ? limit 1 FOR SHARE", checkoutsRows, m.table)
	var resp Checkouts
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, preOrderId)

	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, sqlx.ErrNotFound
	default:
		return nil, err
	}
}
func (m *customCheckoutsModel) CountByUserId(ctx context.Context, userId uint32) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = ?", m.table)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, userId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func (m *customCheckoutsModel) FindByUserId(ctx context.Context, userId uint32, page int32, pageSize int32) ([]*Checkouts, error) {
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`
		SELECT %s FROM %s 
		WHERE user_id = ? 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`, checkoutsRows, m.table)

	var checkouts []*Checkouts
	err := m.conn.QueryRowsCtx(ctx, &checkouts, query, userId, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return checkouts, nil
}
