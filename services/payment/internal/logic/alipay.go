package logic

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"jijizhazha1024/go-mall/services/payment/internal/svc"
	"time"
)

func GenerateAlipayPaymentURL(svcCtx *svc.ServiceContext, PaidAmount int64, expireSeconds int64, pre_order_id string) (string, error) {
	client := svcCtx.Alipay
	// 3. 构造支付订单参数
	var p alipay.TradePagePay
	p.NotifyURL = svcCtx.Config.Alipay.NotifyURL
	p.ReturnURL = svcCtx.Config.Alipay.ReturnURL
	p.Subject = "订单支付" // 可根据实际业务动态设置
	p.OutTradeNo = pre_order_id
	// 将金额从分转换为元，并格式化为字符串（保留两位小数）
	amountYuan := float64(PaidAmount) / 100.0
	p.TotalAmount = fmt.Sprintf("%.2f", amountYuan)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
	// 设置支付订单超时时间，转换为分钟单位的字符串，如 "30m"
	//p.TimeoutExpress = fmt.Sprintf("%dm", expireSeconds/60)
	//fmt.Println(p.TimeoutExpress)
	expireTime := time.Now().Add(time.Duration(expireSeconds) * time.Second)
	p.TimeExpire = expireTime.Format("2006-01-02 15:04:05") // 4. 调用支付宝生成支付订单，返回签名后的订单字符串
	orderString, err := client.TradePagePay(p)
	if err != nil {
		return "", fmt.Errorf("failed to create alipay order: %v", err)
	}
	// 5. 返回生成的支付 URL（订单字符串）
	return orderString.String(), nil
	return "", nil
}
