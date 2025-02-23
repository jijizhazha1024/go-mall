package logic

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	paymentM "jijizhazha1024/go-mall/dal/model/payment"
	"jijizhazha1024/go-mall/services/payment/internal/svc"
	"jijizhazha1024/go-mall/services/payment/payment"
	"time"
)

type CreatePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// convertModelToPaymentItem 将 DAL 模型转换为 proto 定义的 PaymentItem
func convertModelToPaymentItem(p *paymentM.Payments) *payment.PaymentItem {
	var method payment.PaymentMethod
	switch p.PaymentMethod {
	case "alipay":
		method = payment.PaymentMethod_ALIPAY
	case "wx_pay":
		method = payment.PaymentMethod_WECHAT_PAY
	default:
		method = payment.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
	return &payment.PaymentItem{
		PaymentId:      p.PaymentId,
		PreOrderId:     p.PreOrderId,
		OrderId:        p.OrderId.String,
		OriginalAmount: p.OriginalAmount,
		PaidAmount:     p.PaidAmount.Int64,
		PaymentMethod:  method,
		TransactionId:  p.TransactionId.String,
		PayUrl:         p.PayUrl,
		ExpireTime:     p.ExpireTime,
		Status:         payment.PaymentStatus(p.Status),
		CreatedAt:      p.CreatedAt.UnixMilli(),
		UpdatedAt:      p.UpdatedAt.UnixMilli(),
		PaidAt:         p.PaidAt.Int64,
	}
}
func (l *CreatePaymentLogic) CreatePayment(in *payment.PaymentReq) (*payment.PaymentResp, error) {
	// 1. 幂等性校验：根据 idempotency_key 查询是否已经创建过支付单
	existingPayment, err := l.svcCtx.PaymentModel.FindOneByIdempotencyKey(l.ctx, in.IdempotencyKey)
	if err == nil && existingPayment != nil {
		return &payment.PaymentResp{
			StatusCode: 0,
			StatusMsg:  "重复请求，返回已存在的支付单",
			Payment:    convertModelToPaymentItem(existingPayment),
		}, nil
	}

	// 2. 获取预订单信息（实际项目中需查询预订单表或调用预订单服务）
	// 这里简单模拟，假设预订单金额为 10000 分
	//TODO
	originalAmount := int64(10000)
	paidAmount := int64(93223)
	// 3. 生成支付单信息
	paymentId := generateUUID()
	now := time.Now().Unix()
	expireTime := now + 1800 // 支付链接 30 分钟后过期

	// 4. 调用第三方支付生成支付链接（此处根据不同渠道简单模拟返回 URL）
	var payUrl string
	switch in.PaymentMethod {
	case payment.PaymentMethod_WECHAT_PAY:
		payUrl = "https://wechat.pay/" + paymentId
	case payment.PaymentMethod_ALIPAY:
		payUrl, err = GenerateAlipayPaymentURL(l.svcCtx.Config, false, paymentId, paidAmount, 1800)
		if err != nil {
			return nil, err
		}
	default:
		payUrl = "https://default.pay/" + paymentId
	}

	// 5. 构造支付单记录
	newPayment := &paymentM.Payments{

		PaymentId:      paymentId,
		PreOrderId:     in.PreOrderId,
		OrderId:        sql.NullString{}, // 支付成功后更新
		OriginalAmount: originalAmount,
		PaidAmount:     sql.NullInt64{Int64: paidAmount},
		PaymentMethod:  paymentMethodToString(in.PaymentMethod),
		TransactionId:  sql.NullString{},
		PayUrl:         payUrl,
		ExpireTime:     expireTime,
		Status:         int64(payment.PaymentStatus_PAYMENT_STATUS_UNPAID), // 初始状态：待支付
		IdempotencyKey: in.IdempotencyKey,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		PaidAt:         sql.NullInt64{},
	}

	// 6. 插入支付单记录（goctl 生成的 DAL Insert 方法）
	if _, err := l.svcCtx.PaymentModel.Insert(l.ctx, newPayment); err != nil {
		return nil, err
	}

	// 7. 返回创建成功的支付信息
	return &payment.PaymentResp{
		StatusCode: 0,
		StatusMsg:  "支付单创建成功",
		Payment:    convertModelToPaymentItem(newPayment),
	}, nil
}

// paymentMethodToString 将 proto 枚举转换为数据库存储的字符串
func paymentMethodToString(method payment.PaymentMethod) string {
	switch method {
	case payment.PaymentMethod_WECHAT_PAY:
		return "wx_pay"
	case payment.PaymentMethod_ALIPAY:
		return "alipay"
	default:
		return "unknown"
	}
}

// generateUUID 生成一个支付单ID（UUID格式）
func generateUUID() string {
	return uuid.New().String()
}
