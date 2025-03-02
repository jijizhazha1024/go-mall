package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/common/consts/code"
	paymentM "jijizhazha1024/go-mall/dal/model/payment"
	"jijizhazha1024/go-mall/services/order/order"
	"jijizhazha1024/go-mall/services/payment/internal/svc"
	"jijizhazha1024/go-mall/services/payment/payment"
	"net/http"
	"time"
)

type PaymentService struct {
	ctx *svc.ServiceContext
}

func NewPaymentService(ctx *svc.ServiceContext) *PaymentService {
	return &PaymentService{ctx: ctx}
}

// 封装支付宝通知处理
func (s *PaymentService) handleAlipayNotification(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		logx.Infow("Failed to parse form", logx.Field("err", err))
		return
	}
	// DecodeNotification 内部已调用 VerifySign 方法验证签名
	var notify, err = s.ctx.Alipay.DecodeNotification(request.Form)
	if err != nil {
		logx.Errorw("Failed to decode notification", logx.Field("err", err))
		return
	}
	// 根据通知状态处理业务逻辑
	switch notify.TradeStatus {
	case "TRADE_FINISHED":
	// 交易完成（不可退款）
	case "TRADE_CLOSED":
		// 未付款超时关闭/全额退款
		logx.Infow("Payment closed", logx.Field("order_id", notify.OutTradeNo))
	case "TRADE_SUCCESS":
		// 解析时间字符串
		paymentTime, err := time.Parse(time.DateTime, notify.GmtPayment)
		if err != nil {
			logx.Errorw("Failed to parse time", logx.Field("err", err))
			return
		}
		var paymentRes *paymentM.Payments
		timestamp := paymentTime.Unix()
		if err := s.ctx.Model.TransactCtx(request.Context(), func(ctx context.Context, session sqlx.Session) error {
			paymentsModel := s.ctx.PaymentModel.WithSession(session)
			pRes, err := paymentsModel.FindOneByOrderId(ctx, notify.OutTradeNo)
			paymentRes = pRes
			if err != nil {
				logx.Errorw("Failed to find payment record", logx.Field("err", err))
				return err
			}
			switch payment.PaymentStatus(pRes.Status) {
			// 订单状态为待支付时，更新订单状态为已支付，退款
			case payment.PaymentStatus_PAYMENT_STATUS_EXPIRED:
			case payment.PaymentStatus_PAYMENT_STATUS_UNPAID:
				// 支付成功
				if err := paymentsModel.UpdateInfoByOrderId(ctx, &paymentM.Payments{
					OrderId:       sql.NullString{String: notify.OutTradeNo, Valid: true}, // 支付成功后更新
					TransactionId: sql.NullString{String: notify.TradeNo, Valid: true},
					Status:        int64(payment.PaymentStatus_PAYMENT_STATUS_PAID),
					PaidAt:        sql.NullInt64{Int64: timestamp},
				}); err != nil {
					return err
				}
				//状态异常，退款操作
			}
			return nil
		}); err != nil {
			logx.Errorw("Failed to update payment record", logx.Field("err", err), logx.Field("order_id", notify.OutTradeNo))
			return
		}

		orderRes, err := s.ctx.OrderRpc.UpdateOrder2PaymentSuccess(request.Context(), &order.UpdateOrder2PaymentSuccessRequest{
			OrderId: notify.OutTradeNo,
			PaymentResult: &order.PaymentResult{
				TransactionId: notify.TradeNo,
				PaidAmount:    paymentRes.PaidAmount.Int64,
				PaidAt:        timestamp,
			},
		})
		if err != nil {
			logx.Errorw("Failed to update order status", logx.Field("err", err))
			return
		}
		if orderRes.StatusCode != code.Success {
			logx.Errorw("Failed to update order status", logx.Field("err", err))
			return
		}
		logx.Infow("Payment success", logx.Field("order_id", notify.OutTradeNo))

	}
	// 返回确认响应给支付宝
	alipay.ACKNotification(writer)

}

// 封装HTTP服务启动
func (s *PaymentService) startHTTPServer() {
	http.HandleFunc(s.ctx.Config.Alipay.NotifyPath, s.handleAlipayNotification)
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", s.ctx.Config.Alipay.NotifyPort), nil); err != nil {
			logx.Errorw("http server error", logx.Field("err", err))
		}
	}()
}
