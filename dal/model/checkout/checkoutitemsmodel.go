package checkout

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CheckoutItemsModel = (*customCheckoutItemsModel)(nil)

type (
	// CheckoutItemsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCheckoutItemsModel.
	CheckoutItemsModel interface {
		checkoutItemsModel
		withSession(session sqlx.Session) CheckoutItemsModel
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
