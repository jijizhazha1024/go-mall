package coupons

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

// --------------- 获取优惠券 ---------------
func Test_CouponListHandler(t *testing.T) {
	url := "http://127.0.0.1:8009/douyin/coupon/list?page=1&size=10"

	req, _ := http.NewRequest("GET", url, nil)
	// 生成token到test/rpc/auths_test.go/TestAuthenticationLogic_GenerateToken方法
	req.Header.Add("access_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJjbGllbnRfaXAiOiIxMjcuMC4wLjEiLCJleHAiOjE3Mzk3MDM3ODksImlhdCI6MTczOTY5NjU4OX0.BvZh6RBDRkqCNFW2HcFeLPaLD74F5kSvraXupLlQIGY")
	req.Header.Add("refresh_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJjbGllbnRfaXAiOiIxMjcuMC4wLjEiLCJleHAiOjE3NDAyODY3MzgsImlhdCI6MTczOTY4MTkzOH0.RoVN1IXqaGLZUFXhMRMbNlbVy4OEmbvl2YxckVIHbA4")
	req.Header.Add("X-Forward-For", "127.0.0.1")

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Log(err)
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

// --------------- 获取优惠券详细 ---------------

// 优惠券不存在情况
func Test_GetCouponHandler_NotFound(t *testing.T) {
	url := "http://127.0.0.1:8009/douyin/coupon/detail?id=11"

	req, _ := http.NewRequest("GET", url, nil)
	// 生成token到test/rpc/auths_test.go/TestAuthenticationLogic_GenerateToken方法
	req.Header.Add("access_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJjbGllbnRfaXAiOiIxMjcuMC4wLjEiLCJleHAiOjE3Mzk3MDM3ODksImlhdCI6MTczOTY5NjU4OX0.BvZh6RBDRkqCNFW2HcFeLPaLD74F5kSvraXupLlQIGY")
	req.Header.Add("refresh_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJjbGllbnRfaXAiOiIxMjcuMC4wLjEiLCJleHAiOjE3NDAyODY3MzgsImlhdCI6MTczOTY4MTkzOH0.RoVN1IXqaGLZUFXhMRMbNlbVy4OEmbvl2YxckVIHbA4")
	req.Header.Add("X-Forward-For", "127.0.0.1")

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Log(err)
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

// 优惠券存在情况
func Test_GetCouponHandler(t *testing.T) {
	url := "http://127.0.0.1:8009/douyin/coupon/detail?id=FJ20250214001"

	req, _ := http.NewRequest("GET", url, nil)
	// 生成token到test/rpc/auths_test.go/TestAuthenticationLogic_GenerateToken方法
	req.Header.Add("access_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJjbGllbnRfaXAiOiIxMjcuMC4wLjEiLCJleHAiOjE3Mzk3MDM3ODksImlhdCI6MTczOTY5NjU4OX0.BvZh6RBDRkqCNFW2HcFeLPaLD74F5kSvraXupLlQIGY")
	req.Header.Add("refresh_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJjbGllbnRfaXAiOiIxMjcuMC4wLjEiLCJleHAiOjE3NDAyODY3MzgsImlhdCI6MTczOTY4MTkzOH0.RoVN1IXqaGLZUFXhMRMbNlbVy4OEmbvl2YxckVIHbA4")
	req.Header.Add("X-Forward-For", "127.0.0.1")

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Log(err)
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

// --------------- 领取优惠券 ---------------
func Test_ClaimCouponHandler(t *testing.T) {
	url := "http://127.0.0.1:8009/douyin/coupon/claim"
	jsonData, _ := json.Marshal(map[string]string{
		"id": "xxxx",
	})
	// 或使用 NewRequest：
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("access_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJjbGllbnRfaXAiOiIxMjcuMC4wLjEiLCJleHAiOjE3Mzk3MTQ0NzYsImlhdCI6MTczOTcwNzI3Nn0.R-cnBDfuqgXfVkHra3rdvRvst68En4ufRj8TfIsITbw")
	req.Header.Add("refresh_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3QiLCJjbGllbnRfaXAiOiIxMjcuMC4wLjEiLCJleHAiOjE3NDAzMTIwNzYsImlhdCI6MTczOTcwNzI3Nn0.xySgHTw9DWUkMxhBg3rj-wUUCy0BKVw_aAOtDogL3tU")
	req.Header.Add("X-Forward-For", "127.0.0.1")
	res, _ := http.DefaultClient.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Log(err)
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
