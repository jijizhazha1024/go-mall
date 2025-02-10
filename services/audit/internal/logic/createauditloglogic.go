package logic

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/common/utils/metadatactx"
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
	res := &audit.CreateAuditLogRes{} // 简化变量声明
	clientIP := in.GetClientIp()
	if clientIP == "" {
		var ok bool
		clientIP, ok = metadatactx.ExtractFromMetadataCtx(l.ctx, biz.ClientIPKey)
		if !ok {
			res.StatusCode = code.NotWithClientIP
			res.StatusMsg = code.NotWithClientIPMsg
			l.Logger.Infow("client ip is empty", logx.Field("user_id", in.UserId)) // 优化日志记录
			return res, nil
		}
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
		ClientIP:    clientIP,
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
