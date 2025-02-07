package vars

import "errors"

var (
	ErrUnknownCommand = errors.New("未知命令")
	ErrASTParseFailed = errors.New("ast解析失败")
)
