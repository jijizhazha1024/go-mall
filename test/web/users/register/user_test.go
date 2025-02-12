package register

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
	url := "http://localhost:8000/douyin/user/register"
	payload := map[string]string{
		"email":            "user@example.com",
		"password":         "password123",
		"confirm_password": "password123",
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
