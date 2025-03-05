package code

const (
	ProductNotFound = 30000 + iota

	ProductInfoRetrievalFailed
)

const (
	ProductNotFoundMsg            = "商品不存在"
	ProductInfoRetrievalFailedMsg = "商品信息查询失败"
)
