package model

// 通用请求结构（用户原始输入）
type AIRequest struct {
	UserID  string  `json:"user_id"` // 用户身份标识
	Input   string  `json:"input"`   // 自然语言输入（如“我要下单商品A”）
	Context Context `json:"context"` // 请求上下文（可选）
}

type Context struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

// 抽象语法树（AST）接口
type AST interface {
	GetCommandType() string
}

// 下单指令的AST
type OrderAST struct {
	Command   string `json:"command"`    // 指令类型：order
	ProductID string `json:"product_id"` // 商品ID
	Quantity  int    `json:"quantity"`   // 购买数量
	UserID    string `json:"user_id"`    // 用户ID（从请求中继承）
}

func (a *OrderAST) GetCommandType() string { return "order" }

// 查询商品指令的AST
type ProductQueryAST struct {
	Command  string `json:"command"`   // 指令类型：product_query
	Keyword  string `json:"keyword"`   // 搜索关键词
	Page     int    `json:"page"`      // 分页页码
	PageSize int    `json:"page_size"` // 分页大小
	UserID   int    `json:"user_id"`   // 用户ID
}

func (a *ProductQueryAST) GetCommandType() string { return "product_query" }
