package checkout

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CheckoutItemsModel = (*customCheckoutItemsModel)(nil)

type (
	// CheckoutItemsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCheckoutItemsModel.
	CheckoutItemsModel interface {
		checkoutItemsModel
		withSession(session sqlx.Session) CheckoutItemsModel
		FindItemsByUserAndPreOrder(ctx context.Context, userId int32, preOrderId string) ([]CheckoutItems, error)
	}

	customCheckoutItemsModel struct {
		*defaultCheckoutItemsModel
	}
)

// NewCheckoutItemsModel returns a model for the database table.
func NewCheckoutItemsModel(conn sqlx.SqlConn) CheckoutItemsModel {
	return &customCheckoutItemsModel{
		defaultCheckoutItemsModel: newCheckoutItemsModel(conn),
	}
}

func (m *customCheckoutItemsModel) withSession(session sqlx.Session) CheckoutItemsModel {
	return NewCheckoutItemsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *defaultCheckoutItemsModel) FindItemsByUserAndPreOrder(ctx context.Context, userId int32, preOrderId string) ([]CheckoutItems, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `pre_order_id` = ?", checkoutItemsRows, m.table)
	var resp []CheckoutItems
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, preOrderId)
	switch err {
	case nil:
		// Return the found checkout record
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, sqlx.ErrNotFound
	default:
		// If there is another error, return it
		return nil, err
	}
}
