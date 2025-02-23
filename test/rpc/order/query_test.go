package order

import (
	"context"
	"jijizhazha1024/go-mall/common/consts/code"
	"jijizhazha1024/go-mall/services/order/order"
	"testing"
)

func TestListOrders(t *testing.T) {
	t.Run("获取商品列表", func(t *testing.T) {
		listOrders, err := orderClient.ListOrders(context.Background(), &order.ListOrdersRequest{
			Pagination: &order.ListOrdersRequest_Pagination{
				Page:     1,
				PageSize: 10,
			},
			UserId: 1,
		})
		if err != nil {
			t.Error(err)
		}
		if listOrders.StatusCode != code.Success {
			t.Log(listOrders.StatusMsg)
			return
		}
		for _, o := range listOrders.Orders {
			t.Log(o)
		}
	})
	t.Run("获取商品列表_空", func(t *testing.T) {
		// 测试空数据
		listOrders, err := orderClient.ListOrders(context.Background(), &order.ListOrdersRequest{
			Pagination: &order.ListOrdersRequest_Pagination{
				Page:     100,
				PageSize: 10,
			},
			UserId: 1,
		})
		if err != nil {
			t.Error(err)
		}
		if listOrders.StatusCode != code.Success {
			t.Log(listOrders.StatusMsg)
			return
		}
		for _, o := range listOrders.Orders {
			t.Log(o)
		}
	})

}
func TestGetOrder(t *testing.T) {
	t.Run("获取商品详情", func(t *testing.T) {
		orderDetail, err := orderClient.GetOrder(context.Background(), &order.GetOrderRequest{
			OrderId: "1",
			UserId:  1,
		})
		if err != nil {
			t.Error(err)
		}
		if orderDetail.StatusCode != code.Success {
			t.Log(orderDetail.StatusMsg)
			return
		}
		t.Logf("orderDetail: %+v", orderDetail.Order)
		t.Logf("items: %+v", orderDetail.Items)
		t.Logf("addres: %+v", orderDetail.Address)
	})
	t.Run("订单不存在", func(t *testing.T) {
		// 测试空数据
		orderDetail, err := orderClient.GetOrder(context.Background(), &order.GetOrderRequest{
			OrderId: "100",
			UserId:  1,
		})
		if err != nil {
			t.Error(err)
		}
		if orderDetail.StatusCode != code.OrderNotExist {
			t.Log(orderDetail.StatusMsg)
			return
		}
	})
	t.Run("订单内部接口调用", func(t *testing.T) {
		orderDetail, err := orderClient.GetOrder2Payment(context.Background(), &order.GetOrderRequest{
			OrderId: "1",
			UserId:  1,
		})
		if err != nil {
			t.Error(err)
		}
		if orderDetail.StatusCode != code.Success {
			t.Log(orderDetail.StatusMsg)
			return
		}
		t.Logf("orderDetail: %+v", orderDetail.Order)
	})

} .
