package code

const (
	UpdateInventoryError int32 = 100000 + iota
	InventoryNotEnough
	InventoryDecreaseFailed
	ProductNotFoundInventory
	OrderhasBeenPaid
	InvalidParams
)
const (
	UpdateInventoryErrorMsg     = "更新库存失败"
	InventoryNotEnoughMsg       = "库存不足"
	InventoryDecreaseFailedMsg  = "库存减少失败"
	ProductNotFoundInventoryMsg = "商品不存在"
	OrderhasBeenPaidMsg         = "订单已处理"
	InvalidParamsMsg            = "参数错误"
)
