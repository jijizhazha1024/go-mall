package product_query

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"jijizhazha1024/go-mall/common/consts/biz"
	"jijizhazha1024/go-mall/services/ai/internal/core/model"
	"jijizhazha1024/go-mall/services/ai/internal/utils/gpt"
	"os"
	"testing"
)

func TestPrompt(t *testing.T) {

	newGpt := gpt.NewGpt(os.Getenv(biz.ApiKey), os.Getenv(biz.ModelID))
	productQueries := []string{
		"我正在寻找最新的运动鞋，最好是耐克品牌的，价格在300到500元之间，而且要是热门款式的。",
		"想了解一下新款的笔记本电脑，预算大概是4000到6000元，对品牌没有特别要求，但必须是轻薄型的。",
		"搜索适合户外活动使用的双肩背包，期望价位在100到200元，颜色最好是蓝色或黑色。",
		"寻找一款新的智能手表，希望它具备心率监测功能，并且价格不要超过300元。",
		"我对摄影感兴趣，计划购买一台入门级单反相机，价格范围在2000至3000元之间。",
		"想要买一些新出的厨房小电器，例如电热水壶，价格大概在100元左右，希望能找到性价比高的产品。",
		"考虑购入一款新的智能手机，偏好华为品牌，预算为2000到3000元，关注点在于拍照效果要好。",
		"寻找最新款的游戏耳机，价格在300到500元之间，主要用途是玩电脑游戏时使用。",
		"想买一条新的牛仔裤，价格区间为150到250元，重点是舒适度和款式要时尚。",
		"计划购置一套新的床上用品，包括被套和枕套，希望是纯棉材质的，预算大概在300到500元之间。",
	}
	for _, productQuery := range productQueries {
		productQueryAST := model.ProductQueryAST{
			BaseAST: model.BaseAST{
				Command: productQuery,
				UserID:  1,
			},
		}
		response, err := newGpt.ChatWithModel(context.Background(), Prompt, fmt.Sprintf("用户输入：%s", productQuery))
		if err != nil {
			t.Errorf("ChatWithModel error: %v", err)
		}
		t.Logf("response: %s", response)
		if err := jsonx.Unmarshal([]byte(response), &productQueryAST); err != nil {
			t.Errorf("Unmarshal error: %v", err)
		}
		t.Logf("productQueryAST: %+v", productQueryAST)

	}

}
