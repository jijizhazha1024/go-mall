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
		GetOrderByOrderIDAndUserIDWithLock(ctx context.Context, orderId string, userId int32) (*Orders, error)
		GetOrdersByUserID(ctx context.Context, userId, page, size int32) ([]*Orders, error)
		UpdateOrder2Payment(context.Context, string, int32, *order.PaymentResult, order.OrderStatus, order.PaymentStatus) error
		UpdateOrder2PaymentRollback(context.Context, string, int32) error
		UpdateOrderStatusByOrderIDAndUserID(context.Context, string, int32, order.OrderStatus, order.PaymentStatus) error
		CheckOrderExistByPreOrderId(context.Context, string, int32) (bool, error)
		CancelOrder(ctx context.Context, userId int32, orderId, reason string) error
	}

	customOrdersModel struct {
		*defaultOrdersModel
	}
)

func (m *customOrdersModel) CancelOrder(ctx context.Context, userId int32, orderId, reason string) error {
	query := fmt.Sprintf("update %s set `order_status` = ?,`payment_status`,`reason` = ? where `order_id` = ? and `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, order.OrderStatus_ORDER_STATUS_CANCELLED,
		order.PaymentStatus_PAYMENT_STATUS_EXPIRED, reason, orderId, userId)
	return err
}

func (m *customOrdersModel) GetOrdersByUserID(ctx context.Context, userId int32, page, size int32) ([]*Orders, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", ordersRows, m.table)
	var resp []*Orders
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customOrdersModel) CheckOrderExistByPreOrderId(ctx context.Context, preOrderId string, userID int32) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where `pre_order_id` = ? and `user_id` = ? limit 1 for share", m.table)
	var cnt int8
	err := m.conn.QueryRowCtx(ctx, &cnt, query, preOrderId, userID)
	switch {
	case err == nil:
		return cnt > 0, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return false, nil
	default:
		return false, err
	}

}

func (m *customOrdersModel) UpdateOrder2Payment(ctx context.Context, orderID string, userId int32,
	paymentResult *order.PaymentResult, orderStatus order.OrderStatus, paymentStatus order.PaymentStatus) error {
	query := fmt.Sprintf("update %s set `order_status` = ? , `payment_status` = ? ,`transaction_id` = ?, `paid_amount` = ?, `paid_at` = ? where `order_id` = ? and `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, orderStatus,
		paymentStatus,
		paymentResult.TransactionId, paymentResult.PaidAmount,
		paymentResult.PaidAt, orderID, userId)
	return err
}
func (m *customOrdersModel) UpdateOrder2PaymentRollback(ctx context.Context, orderID string, userId int32) error {
	query := fmt.Sprintf("update %s set `order_status` = ? , `payment_status` = ? where `order_id` = ? and `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, order.OrderStatus_ORDER_STATUS_PENDING_PAYMENT, order.PaymentStatus_PAYMENT_STATUS_PAYING, orderID, userId)
	return err
}
func (m *customOrdersModel) GetOrderByOrderIDAndUserIDWithLock(ctx context.Context, orderId string, userId int32) (*Orders, error) {
	var resp Orders
	query := fmt.Sprintf("select %s from %s where `order_id` = ? and `user_id` = ? LIMIT 1 FOR SHARE ",
		ordersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, orderId, userId)
	return &resp, err
}

func (m *customOrdersModel) UpdateOrderStatusByOrderIDAndUserID(ctx context.Context, orderId string,
	userId int32, orderStatus order.OrderStatus, paymentStatus order.PaymentStatus) error {
	query := fmt.Sprintf("update %s set `order_status` = ?,`payment_status` = ? where `order_id` = ? and `user_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, orderStatus, paymentStatus, orderId, userId)
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
