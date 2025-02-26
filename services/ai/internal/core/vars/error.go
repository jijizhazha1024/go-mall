package vars

import "errors"

var (
	ErrUnknownCommand = errors.New("未知命令")
	ErrASTParseFailed = errors.New("ast解析失败")
)
var (
	ErrProductQueryFailed = errors.New("商品查询失败")
)
