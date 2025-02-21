package order

import (
	"context"
	"jijizhazha1024/go-mall/services/order/order"
	"testing"
)

func TestUpdate(t *testing.T) {
	t.Run("订单状态更新为支付中与补偿", func(t *testing.T) {
		res, err := orderClient.UpdateOrder2PaymentStatus(context.Background(), &order.UpdateOrder2PaymentRequest{
			OrderId: "1",
			UserId:  1,
		})
		if err != nil {
			t.Error(err)
		} else {
			t.Log(res)
		}
		t.Log(res)
		rollback, err := orderClient.UpdateOrder2PaymentStatusRollback(context.Background(), &order.UpdateOrder2PaymentRequest{
			OrderId: "1",
			UserId:  1,
		})
		if err != nil {
			t.Error(err)
		} else {
			t.Log(rollback)
		}
	})
}
