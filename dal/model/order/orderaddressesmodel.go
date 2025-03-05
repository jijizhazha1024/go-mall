package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderAddressesModel = (*customOrderAddressesModel)(nil)

type (
	// OrderAddressesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderAddressesModel.
	OrderAddressesModel interface {
		orderAddressesModel
		withSession(session sqlx.Session) OrderAddressesModel
		GetOrderAddressByOrderID(ctx context.Context, orderID string) (*OrderAddresses, error)
		DeleteOrderAddressByOrderID(ctx context.Context, session sqlx.Session, orderID string) error
	}

	customOrderAddressesModel struct {
		*defaultOrderAddressesModel
	}
)

func (m *customOrderAddressesModel) DeleteOrderAddressByOrderID(ctx context.Context, session sqlx.Session, orderID string) error {
	query := fmt.Sprintf("delete from %s where `order_id` = ?", m.table)
	_, err := session.ExecCtx(ctx, query, orderID)
	return err
}

func (m *customOrderAddressesModel) GetOrderAddressByOrderID(ctx context.Context, orderID string) (*OrderAddresses, error) {
	query := fmt.Sprintf("select %s from %s where `order_id` = ?", orderAddressesRows, m.table)
	var resp OrderAddresses
	err := m.conn.QueryRowCtx(ctx, &resp, query, orderID)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// NewOrderAddressesModel returns a model for the database table.
func NewOrderAddressesModel(conn sqlx.SqlConn) OrderAddressesModel {
	return &customOrderAddressesModel{
		defaultOrderAddressesModel: newOrderAddressesModel(conn),
	}
}

func (m *customOrderAddressesModel) withSession(session sqlx.Session) OrderAddressesModel {
	return NewOrderAddressesModel(sqlx.NewSqlConnFromSession(session))
}
