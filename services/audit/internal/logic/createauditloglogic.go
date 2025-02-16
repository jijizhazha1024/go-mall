package logic

import (
	"context"
	"errors"
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
func Validate(req *audit.CreateAuditLogReq) error {
	var errs []error

	// 校验必填字段
	if req.GetUserId() == 0 { // uint32 零值校验
		errs = append(errs, errors.New("user_id is required"))
	}
	if req.GetActionType() == "" {
		errs = append(errs, errors.New("action_type is required"))
	}
	if req.GetTargetTable() == "" {
		errs = append(errs, errors.New("target_table is required"))
	}
	if req.GetTargetId() == 0 { // int64 零值校验
		errs = append(errs, errors.New("target_id is required"))
	}
	if req.GetClientIp() == "" {
		errs = append(errs, errors.New("client_ip is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// CreateAuditLog 创建审计日志
func (l *CreateAuditLogLogic) CreateAuditLog(in *audit.CreateAuditLogReq) (*audit.CreateAuditLogRes, error) {
	// --------------- check ---------------
	if err := Validate(in); err != nil {
		return nil, err
	}
	spanContext := trace.SpanContextFromContext(l.ctx)
	traceID := spanContext.TraceID().String()
	spanID := spanContext.SpanID().String()

	req := mq.AuditReq{
		UserID:      in.UserId,
		TargetTable: in.TargetTable,
		TargetID:    in.TargetId,
		ActionType:  in.ActionType,
		ActionDesc:  in.ActionDescription,
		NewData:     in.NewData,
		OldData:     in.OldData,
		SpanID:      spanID,
		TraceID:     traceID,
		ClientIP:    in.ClientIp,
		CreatedAt:   in.CreateAt,
	}
	if err := l.svcCtx.AuditMQ.Product(&req); err != nil {
		l.Logger.Errorw("CreateAuditLogLogic.CreateAuditLog",
			logx.Field("traceID", traceID),
			logx.Field("err", err))
		return nil, err
	}
	return &audit.CreateAuditLogRes{Ok: true}, nil
}
