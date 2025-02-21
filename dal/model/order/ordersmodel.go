package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/services/order/order"
)

var _ OrdersModel = (*customOrdersModel)(nil)

type (
	// OrdersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrdersModel.
	OrdersModel interface {
		ordersModel
		WithSession(session sqlx.Session) OrdersModel
		GetOrderStatusByOrderIDAndUserIDWithLock(ctx context.Context, orderId string, userId int32) (int64, error)
		UpdateOrderStatusByOrderIDAndUserID(ctx context.Context, orderId string, userId int32, payment order.OrderStatus) error
	}

	customOrdersModel struct {
		*defaultOrdersModel
	}
)

func (m *customOrdersModel) UpdateOrderStatusByOrderIDAndUserID(ctx context.Context, orderId string, userId int32, payment order.OrderStatus) error {
	query := fmt.Sprintf("update %s set `order_status` = ? where `order_id` = ? and `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, payment, orderId, userId)
	return err
}

func (m *customOrdersModel) GetOrderStatusByOrderIDAndUserIDWithLock(ctx context.Context, orderId string, userId int32) (int64, error) {
	var orderStatus int64
	query := fmt.Sprintf("select `order_status` from %s where `order_id` = ? and `user_id` = ? LIMIT 1 FOR SHARE ",
		m.table)
	err := m.conn.QueryRowCtx(ctx, &orderStatus, query, orderId, userId)
	switch {
	case err == nil:
		return orderStatus, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, sqlx.ErrNotFound
	default:
		return 0, err
	}
}

// NewOrdersModel returns a model for the database table.
func NewOrdersModel(conn sqlx.SqlConn) OrdersModel {
	return &customOrdersModel{
		defaultOrdersModel: newOrdersModel(conn),
	}
}

func (m *customOrdersModel) WithSession(session sqlx.Session) OrdersModel {
	return NewOrdersModel(sqlx.NewSqlConnFromSession(session))
}
