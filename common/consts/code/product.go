package code

const (
	ProductCreated             = 30001
	ProductCreationFailed      = 30002
	ProductUpdated             = 30003
	ProductUpdateFailed        = 30004
	ProductDeleted             = 30005
	ProductDeletionFailed      = 30006
	ProductInfoRetrieved       = 30007
	ProductInfoRetrievalFailed = 30008
	ProductSensitiveWordFailed = 30009
	EsFailed                   = 30010
	ProductCacheFailed         = 30011
)

const (
	ProductCreatedMsg             = "商品创建成功"
	ProductCreationFailedMsg      = "商品创建失败"
	ProductUpdatedMsg             = "商品信息更新成功"
	ProductUpdateFailedMsg        = "商品信息更新失败"
	ProductDeletedMsg             = "商品删除成功"
	ProductDeletionFailedMsg      = "商品删除失败"
	ProductInfoRetrievedMsg       = "商品信息查询成功"
	ProductInfoRetrievalFailedMsg = "商品信息查询失败"
	ProductSensitiveWordFailedMsg = "敏感词校验失败"
	EsFailedMag                   = "Es操作失败"
	ProductCacheFailedMsg         = "缓存操作失败"
)
