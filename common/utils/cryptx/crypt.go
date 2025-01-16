package cryptx

import (
	"fmt"
	"golang.org/x/crypto/scrypt"
)

// salt 密码盐
const salt = "go-mall"

// PasswordEncrypt 密码加密
func PasswordEncrypt(password string) string {
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	return fmt.Sprintf("%x", string(dk))
}

// PasswordVerify 密码验证
func PasswordVerify(password, hash string) bool {
	dk, _ := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	return fmt.Sprintf("%x", string(dk)) == hash
}
