package model

// 通用响应结构
type AIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`   // 具体响应数据
	Error   string      `json:"error"`  // 错误信息
}

// 下单响应数据
type OrderResponseData struct {
	OrderID     string `json:"order_id"`
	TotalPrice  int    `json:"total_price"`
	Status      string `json:"status"`      // 如 "created"、"paid"
	IsMock      bool   `json:"is_mock"`     // 是否为模拟操作
}

// 商品查询响应数据
type ProductQueryResponseData struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`      // 总商品数
}

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}