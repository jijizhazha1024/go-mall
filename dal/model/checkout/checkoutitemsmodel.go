package checkout

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jijizhazha1024/go-mall/services/checkout/checkout"
	"strings"
)

var _ CheckoutItemsModel = (*customCheckoutItemsModel)(nil)

type (
	// CheckoutItemsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCheckoutItemsModel.
	CheckoutItemsModel interface {
		checkoutItemsModel
		WithSession(session sqlx.Session) CheckoutItemsModel
		FindItemsByPreOrder(ctx context.Context, preOrderId string) ([]*CheckoutItems, error)
		FindItemsByPreOrderIds(ctx context.Context, preOrderIds []string) (map[string][]*checkout.CheckoutItem, error)
	}

	customCheckoutItemsModel struct {
		*defaultCheckoutItemsModel
	}
)

// NewCheckoutItemsModel returns a model for the database table.
func NewCheckoutItemsModel(conn sqlx.SqlConn) CheckoutItemsModel {
	return &customCheckoutItemsModel{
		defaultCheckoutItemsModel: newCheckoutItemsModel(conn),
	}
}

func (m *customCheckoutItemsModel) WithSession(session sqlx.Session) CheckoutItemsModel {
	return NewCheckoutItemsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customCheckoutItemsModel) FindItemsByPreOrder(ctx context.Context, preOrderId string) ([]*CheckoutItems, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE pre_order_id = ?", m.table)
	var items []*CheckoutItems
	err := m.conn.QueryRowsCtx(ctx, &items, query, preOrderId)
	if err != nil {
		return nil, err
	}
	return items, nil
}
func (m *customCheckoutItemsModel) FindItemsByPreOrderIds(ctx context.Context, preOrderIds []string) (map[string][]*checkout.CheckoutItem, error) {
	preOrders := strings.Join(preOrderIds, ",")
	query := fmt.Sprintf("SELECT * FROM %s WHERE pre_order_id IN (?)", m.table)
	items := make([]*CheckoutItems, 0)
	err := m.conn.QueryRowsCtx(ctx, &items, query, preOrders)
	if err != nil {
		return nil, err
	}
	// 组装 map[pre_order_id] -> []CheckoutItem
	itemsMap := make(map[string][]*checkout.CheckoutItem)
	for _, item := range items {
		// 解析 snapshot
		var snapshotData map[string]string
		if err := json.Unmarshal([]byte(item.Snapshot), &snapshotData); err != nil {
			snapshotData = map[string]string{} // 解析失败，设置为空 map
		}

		// 获取 name 和 specs
		productName := snapshotData["name"]
		productDesc := snapshotData["specs"]

		// 组装返回对象
		itemsMap[item.PreOrderId] = append(itemsMap[item.PreOrderId], &checkout.CheckoutItem{
			ProductId:   int32(item.ProductId),
			Quantity:    int32(item.Quantity),
			ProductName: productName,
			ProductDesc: productDesc,
			Price:       item.Price,
		})
	}

	return itemsMap, nil
}
