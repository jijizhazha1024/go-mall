package biz

import "time"

const (
	ProductRpcPort     = 10002
	ProductEsIndexName = "products"
	ProductRedisPVName = "productPV"
	ScanProductPVTime  = 5 * time.Hour
)
const (
	ProductIDKey       = "product:%d"
	ProductIDKeyExpire = 60 * 60 * 2 // 2小时 过期缓存，由回写触发更新

)
