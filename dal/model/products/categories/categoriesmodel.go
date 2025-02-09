package categories

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CategoriesModel = (*customCategoriesModel)(nil)

type (
	// CategoriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCategoriesModel.
	CategoriesModel interface {
		categoriesModel
		withSession(session sqlx.Session) CategoriesModel
	}

	customCategoriesModel struct {
		*defaultCategoriesModel
	}
)

// NewCategoriesModel returns a model for the database table.
func NewCategoriesModel(conn sqlx.SqlConn) CategoriesModel {
	return &customCategoriesModel{
		defaultCategoriesModel: newCategoriesModel(conn),
	}
}

func (m *customCategoriesModel) withSession(session sqlx.Session) CategoriesModel {
	return NewCategoriesModel(sqlx.NewSqlConnFromSession(session))
}
