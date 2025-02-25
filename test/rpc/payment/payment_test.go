package payment

import (
	"context"
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"jijizhazha1024/go-mall/services/payment/payment"
	"testing"
)

var payment_client payment.PaymentClient

func initpayment() {
	conn, err := grpc.NewClient(fmt.Sprintf("0.0.0.0:%d", 10006),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	payment_client = payment.NewPaymentClient(conn)
}
func TestCreatePayment(t *testing.T) {
	initpayment()
	resp, err := payment_client.CreatePayment(context.Background(), &payment.PaymentReq{
		UserId:         02131,
		PreOrderId:     "2312342111110111",
		PaymentMethod:  payment.PaymentMethod_ALIPAY,
		IdempotencyKey: "1234567892111110111",
	})
	if err != nil {
		return
	}
	fmt.Println(resp)
	if err != nil {
		t.Fatal(err)

	}
	fmt.Println(" success", resp)
	t.Log(" success", resp)
}
func TestGetallProduct(t *testing.T) {
	var appid string = "9021000144642851"
	var privatekey string
	var publickey string
	initpayment()
	client, err := alipay.New(appid, privatekey, false)
	if err != nil {
		fmt.Printf("失败1")
		return
	}
	err = client.LoadAliPayPublicKey(publickey)
	if err != nil {
		fmt.Printf("失败2")

		return
	}
	// 3. 构造支付订单参数
	var p alipay.TradePagePay
	p.NotifyURL = "http://gk8fne.natappfree.cc/notify"
	p.Subject = "订单支付" // 可根据实际业务动态设置
	// 将金额从分转换为元，并格式化为字符串（保留两位小数）
	amountYuan := float64(10000) / 100.0
	p.TotalAmount = fmt.Sprintf("%.2f", amountYuan)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY" // 针对移动支付
	p.OutTradeNo = "123456789011"            // 商户订单号，需保证唯一性

	// 设置支付订单超时时间，转换为分钟单位的字符串，如 "30m"
	// 4. 调用支付宝生成支付订单，返回签名后的订单字符串
	orderString, err := client.TradePagePay(p)

	if err != nil {
		fmt.Printf("失败3")

		t.Fatal(err)
		return
	}
	fmt.Println(1)

	fmt.Println(orderString)
	fmt.Println(2)

}
