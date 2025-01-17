package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLogin(t *testing.T) {
	url := "http://localhost:8000/douyin/user/login"
	payload := map[string]string{
		"email":    "user@example.com",
		"password": "password123",
	}
	jsonPayload, err := json.Marshal(payload)
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	assert.Nil(t, err)
	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	fmt.Println(string(body))
}
