package biz

import "errors"

const (
	InventoryRpcPort = 10011
)
const (
	InventoryKeyPrefix        = "inventory:%d"
	InventoryAccessKeyPrefix  = "inventory_access:%d"
	InventoryDeductLockPrefix = "inventory:deduct:lock"
	InventoryProductKey       = "inventory:product"
)

var (
	// InventoryNotEnoughErr 库存不足err
	InventoryNotEnoughErr = errors.New("not enough inventory")
	// InventoryDecreaseFailedErr 扣减失败
	InventoryDecreaseFailedErr = errors.New("decrease inventory failed")
	// InvalidInventoryErr 非法的库存信息
	InvalidInventoryErr = errors.New("invalid inventory")
)
