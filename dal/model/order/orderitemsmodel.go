package order

import (
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
	}

	customOrderItemsModel struct {
		*defaultOrderItemsModel
	}
)

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
