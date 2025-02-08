package audit

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AuditModel = (*customAuditModel)(nil)

type (
	// AuditModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuditModel.
	AuditModel interface {
		auditModel
		withSession(session sqlx.Session) AuditModel
	}

	customAuditModel struct {
		*defaultAuditModel
	}
)

// NewAuditModel returns a model for the database table.
func NewAuditModel(conn sqlx.SqlConn) AuditModel {
	return &customAuditModel{
		defaultAuditModel: newAuditModel(conn),
	}
}

func (m *customAuditModel) withSession(session sqlx.Session) AuditModel {
	return NewAuditModel(sqlx.NewSqlConnFromSession(session))
}
