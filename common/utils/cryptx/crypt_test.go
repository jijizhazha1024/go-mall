package cryptx

import (
	"fmt"
	"testing"
)

func TestPasswordVerify(t *testing.T) {
	for i := 0; i < 10; i++ {
		password := fmt.Sprintf("password%d", i)
		hash := PasswordEncrypt(password)
		if PasswordVerify(password, hash) {
			t.Log("密码验证成功")
		} else {
			t.Error("密码验证失败")
		}
	}
}
