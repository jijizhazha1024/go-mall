package order

import (
	"context"
	"github.com/dtm-labs/client/dtmgrpc"
	_ "github.com/dtm-labs/driver-gozero"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/zrpc"
	"jijizhazha1024/go-mall/services/order/order"
	"testing"
	"time"
)

func TestUpdateOrder(t *testing.T) {
	t.Run("订单状态更新为支付中与补偿", func(t *testing.T) {
		res, err := orderClient.UpdateOrder2PaymentStatus(context.Background(), &order.UpdateOrder2PaymentRequest{
			OrderId: "1",
			UserId:  1,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(res)
		rollback, err := orderClient.UpdateOrder2PaymentStatusRollback(context.Background(), &order.UpdateOrder2PaymentRequest{
			OrderId: "1",
			UserId:  1,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(rollback)

	})
	t.Run("订单状态更新为支付成功与补偿", func(t *testing.T) {
		res, err := orderClient.UpdateOrder2PaymentSuccess(context.Background(), &order.UpdateOrder2PaymentSuccessRequest{
			OrderId: "1",
			UserId:  1,
			PaymentResult: &order.PaymentResult{
				TransactionId: "xxxx",
				PaidAmount:    100,
				PaidAt:        time.Now().Unix(),
			},
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(res)
		rollback, err := orderClient.UpdateOrder2PaymentSuccessRollback(context.Background(), &order.UpdateOrder2PaymentSuccessRequest{
			OrderId: "1",
			UserId:  1,
			PaymentResult: &order.PaymentResult{
				TransactionId: "xxxx",
				PaidAmount:    100,
				PaidAt:        time.Now().Unix(),
			},
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(rollback)
	})
}

func TestDtmSaga(t *testing.T) {
	orderRpc := zrpc.RpcClientConf{
		Target: "consul://124.71.72.124:8500/order.rpc", NonBlock: true,
	}
	target, err := orderRpc.BuildTarget()
	if err != nil {
		t.Error(err)
		return
	}
	dtmServer := "consul://124.71.72.124:8500/dtmservice"
	sagaGrpc := dtmgrpc.NewSagaGrpc(dtmServer, uuid.New().String()).
		Add(target+order.OrderService_UpdateOrder2PaymentStatus_FullMethodName,
			target+order.OrderService_UpdateOrder2PaymentStatusRollback_FullMethodName,
			&order.UpdateOrder2PaymentRequest{UserId: 1, OrderId: "xxxx"})
	sagaGrpc.WaitResult = true
	// eer => status.Error()
	if err := sagaGrpc.Submit(); err != nil {
		return
	}

	t.Log(sagaGrpc.Gid)
}
