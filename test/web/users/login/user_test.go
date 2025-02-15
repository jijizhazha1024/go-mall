package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	url := "http://localhost:8001/douyin/user/login"
	payload := map[string]string{
		"email":    "djj2@example.com",
		"password": "111111",
	}
	jsonPayload, err := json.Marshal(payload)
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	assert.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	fmt.Println(string(body))
	assert.Equal(t, http.StatusOK, res.StatusCode, "预期状态码200")
}
