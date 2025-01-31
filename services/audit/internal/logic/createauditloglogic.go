package logic

import (
	"context"
	"jijizhazha1024/go-mall/services/audit/audit"
	"jijizhazha1024/go-mall/services/audit/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAuditLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAuditLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAuditLogLogic {
	return &CreateAuditLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateAuditLog 创建审计日志
func (l *CreateAuditLogLogic) CreateAuditLog(in *audit.CreateAuditLogReq) (*audit.CreateAuditLogRes, error) {
	// todo: add your logic here and delete this line
	l.Info(in)
	return &audit.CreateAuditLogRes{}, nil
}
