package shopping

import (
	"github.com/shopspring/decimal"
)

// FenToYuan 将分转换为元字符串，保留两位小数（支持负数）
func FenToYuan(fen int64) string {
	// 创建 decimal 类型的分值
	fenDecimal := decimal.NewFromInt(fen)

	// 转换为元并保留两位小数
	// 等价于 fen / 100，但避免浮点精度问题
	yuan := fenDecimal.Div(decimal.NewFromInt(100))

	// 格式化为带两位小数的字符串
	return yuan.StringFixedBank(2) // 银行家舍入法（四舍六入五成双）
}
