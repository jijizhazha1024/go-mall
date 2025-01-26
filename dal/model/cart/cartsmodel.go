package cart

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CartsModel = (*customCartsModel)(nil)

type (
	// CartsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCartsModel.
	CartsModel interface {
		cartsModel
		withSession(session sqlx.Session) CartsModel
	}

	customCartsModel struct {
		*defaultCartsModel
	}
)

// NewCartsModel returns a model for the database table.
func NewCartsModel(conn sqlx.SqlConn) CartsModel {
	return &customCartsModel{
		defaultCartsModel: newCartsModel(conn),
	}
}

func (m *customCartsModel) withSession(session sqlx.Session) CartsModel {
	return NewCartsModel(sqlx.NewSqlConnFromSession(session))
}
