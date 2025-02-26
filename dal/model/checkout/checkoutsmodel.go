package checkout

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CheckoutsModel = (*customCheckoutsModel)(nil)

type (
	// CheckoutsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCheckoutsModel.
	CheckoutsModel interface {
		checkoutsModel
		withSession(session sqlx.Session) CheckoutsModel
		UpdateStatus(ctx context.Context, status int64, preOrderId string) error
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
func (m *customCheckoutsModel) UpdateStatus(ctx context.Context, status int64, preOrderId string) error {
	updateQuery := "UPDATE checkouts SET status = ? WHERE pre_order_id = ?"

	_, err := m.conn.ExecCtx(ctx, updateQuery, status, preOrderId)
	if err != nil {
		return err
	}

	return nil
}
