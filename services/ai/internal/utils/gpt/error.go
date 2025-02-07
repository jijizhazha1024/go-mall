package gpt

import "errors"

var (
	ErrBalanceInsufficient = errors.New("余额不足")
	ErrModelNotFound       = errors.New("模型不存在")
	ErrModelCallFailed     = errors.New("模型调用失败")
	ErrModelCallTimeout    = errors.New("模型调用超时")
)
