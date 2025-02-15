package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	url := "http://localhost:8001/douyin/user/register"

	// 1. 构造符合验证规则的请求体
	payload := map[string]string{
		"email":           "djj2@example.com", // 确保邮箱格式有效
		"password":        "password123",      // 长度>=8
		"confirmPassword": "password123",      // 与password一致
	}
	jsonPayload, err := json.Marshal(payload)
	assert.Nil(t, err)

	// 2. 创建带超时的HTTP客户端
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 3. 设置请求头
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")

	// 4. 发送请求并处理响应
	res, err := client.Do(req)
	if !assert.Nil(t, err) {
		return
	}
	defer res.Body.Close()

	// 5. 解析响应内容
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	fmt.Println("响应内容:", string(body))

	// 6. 验证响应状态码
	assert.Equal(t, http.StatusOK, res.StatusCode, "预期状态码200")
}
