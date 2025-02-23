package audit

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AuditModel = (*customAuditModel)(nil)

type (
	// AuditModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAuditModel.
	AuditModel interface {
		auditModel
		withSession(session sqlx.Session) AuditModel
		CheckExistByTraceID(ctx context.Context, traceID string) (bool, error)
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
func (m *customAuditModel) CheckExistByTraceID(ctx context.Context, traceID string) (bool, error) {
	var cnt int64
	query := fmt.Sprintf("select count(*) from %s where `trace_id` = ?", m.table)
	if err := m.conn.QueryRowCtx(ctx, &cnt, query, traceID); err != nil {
		return false, err
	}
	return cnt > 0, nil
}
