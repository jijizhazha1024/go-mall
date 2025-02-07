package model

type BaseAST struct {
	Command string `json:"command"` // 用户原始输入
	UserID  int    `json:"user_id"` // 用户ID
}

type AST interface {
	GetCommandType() string
}

// ProductQueryAST 查询商品指令的AST
type ProductQueryAST struct {
	BaseAST
	Conditions struct {
		Name     string   `json:"name"`     // 商品名称
		New      bool     `json:"new"`      // 是否新品
		Hot      bool     `json:"hot"`      // 是否热门
		Category []string `json:"category"` // 商品分类
		Price    struct {
			Min int `json:"min"` // 最低价格
			Max int `json:"max"` // 最高价格
		} `json:"price"` // 价格范围
		Keyword string `json:"keyword"` // 搜索关键词 （模糊查询）
	} `json:"conditions"`
}

func (a *ProductQueryAST) GetCommandType() string {
	return "product_query"
}
