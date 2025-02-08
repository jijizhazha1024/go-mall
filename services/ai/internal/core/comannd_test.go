package core

import (
	"context"
	"jijizhazha1024/go-mall/services/ai/internal/config"
	"testing"
)

func TestHandler(t *testing.T) {
	var conf config.Config
	c := NewCommand(&conf)
	handler, err := c.Handler(context.Background(),
		"想了解一下新款的笔记本电脑，预算大概是4000到6000元，对品牌没有特别要求，但必须是轻薄型的。", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(handler)
}
