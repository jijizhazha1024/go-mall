package payment

import (
	"context"
	"fmt"
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
		UserId:        1,
		OrderId:       "199112649961",
		PaymentMethod: payment.PaymentMethod_ALIPAY,
	})

	if err != nil {
		t.Fatal(err)
	}
	t.Log(" success", resp)
}
func TestGetAllProduct(t *testing.T) {
	initpayment()
	listPayments, err := payment_client.ListPayments(context.Background(), &payment.PaymentListReq{
		Pagination: &payment.PaymentListReq_Pagination{
			Page:     1,
			PageSize: 10,
		},
		UserId: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success", listPayments)

}
