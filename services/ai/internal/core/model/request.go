package model

// AIRequest 通用请求结构（用户原始输入）
type AIRequest struct {
	UserID  string `json:"user_id"` // 用户身份标识
	Command string `json:"command"` // 自然语言输入（如“我要下单商品A”）
}
