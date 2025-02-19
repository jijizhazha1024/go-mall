package biz

import "time"

const (
	ProductRpcPort     = 10002
	ProductEsIndexName = "products"
	ProductRedisPVName = "productPV"
	ScanProductPVTime  = 5 * time.Hour
)
