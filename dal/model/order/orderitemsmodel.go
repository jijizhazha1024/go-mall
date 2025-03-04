package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderItemsModel = (*customOrderItemsModel)(nil)

type (
	// OrderItemsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderItemsModel.
	OrderItemsModel interface {
		orderItemsModel
		WithSession(session sqlx.Session) OrderItemsModel
		// BulkInsert 批量插入
		BulkInsert(session sqlx.Session, items []*OrderItems) error
		QueryOrderItemsByOrderID(ctx context.Context, orderID string) ([]*OrderItems, error)
		DeleteOrderItemByOrderID(ctx context.Context, session sqlx.Session, orderID string) error
	}

	customOrderItemsModel struct {
		*defaultOrderItemsModel
	}
)

func (m *customOrderItemsModel) DeleteOrderItemByOrderID(ctx context.Context, session sqlx.Session, orderID string) error {
	query := fmt.Sprintf("delete from %s where `order_id` = ?", m.table)
	_, err := session.ExecCtx(ctx, query, orderID)
	return err
}

func (m *customOrderItemsModel) QueryOrderItemsByOrderID(ctx context.Context, orderID string) ([]*OrderItems, error) {
	query := fmt.Sprintf("select %s from %s where `order_id` = ?", orderItemsRows, m.table)
	var resp []*OrderItems
	err := m.conn.QueryRowsCtx(ctx, &resp, query, orderID)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customOrderItemsModel) BulkInsert(session sqlx.Session, items []*OrderItems) error {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, orderItemsRowsExpectAutoSet)
	bulkInserter, err := sqlx.NewBulkInserter(sqlx.NewSqlConnFromSession(session), query)
	if err != nil {
		return err
	}
	for _, item := range items {
		err = bulkInserter.Insert(item.OrderId, item.ProductId, item.Quantity, item.Price, item.ProductName, item.ProductDesc)
		if err != nil {
			return err
		}
	}
	bulkInserter.Flush()
	return nil
}

// NewOrderItemsModel returns a model for the database table.
func NewOrderItemsModel(conn sqlx.SqlConn) OrderItemsModel {
	return &customOrderItemsModel{
		defaultOrderItemsModel: newOrderItemsModel(conn),
	}
}

func (m *customOrderItemsModel) WithSession(session sqlx.Session) OrderItemsModel {
	return NewOrderItemsModel(sqlx.NewSqlConnFromSession(session))
}
