package order

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ OrderAddressesModel = (*customOrderAddressesModel)(nil)

type (
	// OrderAddressesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderAddressesModel.
	OrderAddressesModel interface {
		orderAddressesModel
		withSession(session sqlx.Session) OrderAddressesModel
	}

	customOrderAddressesModel struct {
		*defaultOrderAddressesModel
	}
)

// NewOrderAddressesModel returns a model for the database table.
func NewOrderAddressesModel(conn sqlx.SqlConn) OrderAddressesModel {
	return &customOrderAddressesModel{
		defaultOrderAddressesModel: newOrderAddressesModel(conn),
	}
}

func (m *customOrderAddressesModel) withSession(session sqlx.Session) OrderAddressesModel {
	return NewOrderAddressesModel(sqlx.NewSqlConnFromSession(session))
}
