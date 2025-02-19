package code

const (
	// 购物车服务

	CartCreated                    = 40001
	CartCreationFailed             = 40002
	CartCleared                    = 40003
	CartClearFailed                = 40004
	CartInfoRetrieved              = 40005
	CartInfoRetrievalFailed        = 40006
	CartSub                        = 40007
	CartSubFailed                  = 40008
	CartProductInfoRetrievalFailed = 40009
	CartProductQuantityInfoFailed  = 40010
	InsufficientInventoryOfProduct = 40011
	CartItemNotFound               = 40012
)
const (
	CartCreatedMsg                    = "购物车商品增加成功"
	CartCreationFailedMsg             = "购物车商品增加失败"
	CartSubMsg                        = "购物车商品扣减成功"
	CartSubFailedMsg                  = "购物车商品扣减失败"
	CartClearedMsg                    = "购物车删除成功"
	CartClearFailedMsg                = "购物车删除失败"
	CartInfoRetrievedMsg              = "购物车信息获取成功"
	CartInfoRetrievalFailedMsg        = "购物车信息获取失败"
	CartProductInfoRetrievalFailedMsg = "购物车商品不存在"
	CartProductQuantityInfoFailedMsg  = "购物车商品数量获取失败"
	InsufficientInventoryOfProductMsg = "商品库存不足"
	CartItemNotFoundMsg               = "购物车商品不存在"
)
