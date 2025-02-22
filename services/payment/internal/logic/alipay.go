package logic

import (
	/*"fmt"
	"github.com/smartwalle/alipay/v3"*/
	"jijizhazha1024/go-mall/services/payment/internal/config"
)

func GenerateAlipayPaymentURL(c config.Config, isProduction bool, paymentId string, PaidAmount int64, expireSeconds int64) (string, error) {
	/*// 1. 创建支付宝客户端
	client, err := alipay.New(c.Alipay.AppId, c.Alipay.PrivateKey, isProduction)
	if err != nil {
		return "", fmt.Errorf("failed to create alipay client: %v", err)
	}
	// 2. 加载支付宝公钥用于验签
	if err = client.LoadAliPayPublicKey(c.Alipay.AlipayPublicKey); err != nil {
		return "", fmt.Errorf("failed to load alipay public key: %v", err)
	}

	// 3. 构造支付订单参数
	var p alipay.TradeAppPay
	p.NotifyURL = c.Alipay.NotifyURL
	p.ReturnURL = c.Alipay.ReturnURL
	p.Subject = "订单支付" // 可根据实际业务动态设置
	p.OutTradeNo = paymentId
	// 将金额从分转换为元，并格式化为字符串（保留两位小数）
	amountYuan := float64(PaidAmount) / 100.0
	p.TotalAmount = fmt.Sprintf("%.2f", amountYuan)
	p.ProductCode = "QUICK_MSECURITY_PAY" // 针对移动支付
	// 设置支付订单超时时间，转换为分钟单位的字符串，如 "30m"
	p.TimeoutExpress = fmt.Sprintf("%dm", expireSeconds/60)
	// 4. 调用支付宝生成支付订单，返回签名后的订单字符串
	orderString, err := client.TradeAppPay(p)
	if err != nil {
		return "", fmt.Errorf("failed to create alipay order: %v", err)
	}

	// 5. 返回生成的支付 URL（订单字符串）
	return orderString, nil*/
	return "", nil
}
