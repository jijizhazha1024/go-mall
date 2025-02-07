package core

import (
	"context"
	"jijizhazha1024/go-mall/services/ai/internal/config"
	"testing"
)

func TestHandler(t *testing.T) {
	var conf config.Config
	conf.Gpt.ApiKey = "5b5ab09c-7298-40d7-b60e-433d21314f36"
	conf.Gpt.ModelID = "ep-20241002090911-md25k"
	c := NewCommand(&conf)
	handler, err := c.Handler(context.Background(),
		"想了解一下新款的笔记本电脑，预算大概是4000到6000元，对品牌没有特别要求，但必须是轻薄型的。", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(handler)
}
