package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	paymentM "jijizhazha1024/go-mall/dal/model/payment"
	"jijizhazha1024/go-mall/services/order/order"
	"jijizhazha1024/go-mall/services/payment/internal/mq"
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
func ConvertModelToPaymentItem(p *paymentM.Payments) *payment.PaymentItem {
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
	// 1. 幂等性校验：根据 OrderId 查询是否已经创建过支付单
	// 2. 锁
	res := &payment.PaymentResp{}
	existingPayment, err := l.svcCtx.PaymentModel.FindOneByOrderId(l.ctx, in.OrderId)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		l.Logger.Errorw("find one by order id failed", logx.Field("err", err))
		return nil, err

	}
	if !errors.Is(err, sqlx.ErrNotFound) {
		res.StatusCode = code.PaymentExist
		res.StatusMsg = code.PaymentExistMsg
		res.Payment = ConvertModelToPaymentItem(existingPayment)
		// 幂等性校验通过，直接返回已存在的支付单
		return res, nil
	}
	// 2. 获取预订单信息（调用订单服务）
	getOrderInfo, err := l.svcCtx.OrderRpc.GetOrder2Payment(l.ctx, &order.GetOrderRequest{
		OrderId: in.OrderId,
		UserId:  in.UserId,
	})
	if err != nil {
		l.Logger.Errorw("get order info failed", logx.Field("err", err))
		return nil, err
	}
	if getOrderInfo.StatusCode != code.Success {
		res.StatusCode = getOrderInfo.StatusCode
		res.StatusMsg = getOrderInfo.StatusMsg
		return res, nil
	}

	originalAmount := getOrderInfo.Order.OriginalAmount
	payableAmount := getOrderInfo.Order.PayableAmount
	fmt.Println(payableAmount)
	// 3. 生成支付单信息
	paymentId := generateUUID()
	// 4. 调用第三方支付生成支付链接（此处根据不同渠道简单模拟返回 URL）
	if in.PaymentMethod != payment.PaymentMethod_ALIPAY {
		res.StatusCode = code.PaymentMethodNotSupport
		res.StatusMsg = code.PaymentMethodNotSupportMsg
	}
	payUrl, err := GenerateAlipayPaymentURL(l.svcCtx, payableAmount, 1800, in.OrderId)
	if err != nil {
		return nil, err
	}
	// 5. 构造支付单记录
	newPayment := &paymentM.Payments{
		UserId:         uint64(in.UserId),
		PaymentId:      paymentId,
		PreOrderId:     getOrderInfo.Order.PreOrderId,
		OrderId:        sql.NullString{String: in.OrderId, Valid: true}, // 支付成功后更新
		OriginalAmount: originalAmount,
		PaidAmount:     sql.NullInt64{Int64: payableAmount, Valid: true},
		PaymentMethod:  PaymentMethodToString(in.PaymentMethod),
		PayUrl:         payUrl,
		ExpireTime:     time.Now().Add(30 * time.Minute).Unix(),
		Status:         int64(payment.PaymentStatus_PAYMENT_STATUS_UNPAID), // 初始状态：待支付
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if _, err := l.svcCtx.PaymentModel.Insert(l.ctx, newPayment); err != nil {
		return nil, err
	}

	if err := l.svcCtx.PaymentMQ.Product(&mq.PaymentReq{
		OrderId: in.OrderId,
	}); err != nil {
		l.Logger.Errorw("product payment failed", logx.Field("err", err))
		return nil, err
	}

	// 7. 返回创建成功的支付信息
	return &payment.PaymentResp{
		Payment: ConvertModelToPaymentItem(newPayment),
	}, nil
}

// PaymentMethodToString paymentMethodToString 将 proto 枚举转换为数据库存储的字符串
func PaymentMethodToString(method payment.PaymentMethod) string {
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
	var uid uuid.UUID
	uid, err := uuid.NewV7()
	if err != nil {
		logx.Errorw("uuid generate failed", logx.Field("err", err))
		uid = uuid.New()
	}
	return uid.String()
}
