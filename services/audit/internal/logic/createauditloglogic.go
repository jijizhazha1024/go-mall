package logic

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"jijizhazha1024/go-mall/services/audit/audit"
	"jijizhazha1024/go-mall/services/audit/internal/mq"
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
	// 1. 获取ClientIP
	clientIP := ""
	spanContext := trace.SpanContextFromContext(l.ctx)
	traceID := spanContext.TraceID().String()
	spanID := spanContext.SpanID().String()

	req := mq.AuditReq{
		UserID:      in.UserId,
		UserName:    in.Username,
		TargetTable: in.TargetTable,
		TargetID:    in.TargetId,
		ActionType:  in.ActionType,
		ActionDesc:  in.ActionDescription,
		NewData:     in.NewData,
		OldData:     in.OldData,
		SpanID:      spanID,
		TraceID:     traceID,
		ClientIP:    clientIP,
		CreatedAt:   in.CreateAt,
	}
	if err := l.svcCtx.AuditMQ.Product(&req); err != nil {
		logx.Errorw("CreateAuditLogLogic.CreateAuditLog",
			logx.Field("traceID", l.ctx.Value("traceID")),
			logx.Field("err", err.Error()))
		return nil, err
	}
	return &audit.CreateAuditLogRes{Ok: true}, nil
}
