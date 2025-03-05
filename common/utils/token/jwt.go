package token

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// Claims 定义一个结构体用于保存负载（用户信息）
type Claims struct {
	UserID   uint32 `json:"user_id"`
	UserName string `json:"username"`
	ClientIP string `json:"client_ip"`
	jwt.RegisteredClaims
}

const secretKey = "go-mall"

func GenerateJWT(userID uint32, userName, clientIP string, expire time.Duration) (string, error) {
	// 设置JWT的负载（Claims）
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		ClientIP: clientIP,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// 创建一个新的Token对象，指定签名方法和Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥对Token进行签名
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ParseJWT(tokenString string) (*Claims, error) {
	// 解析并验证 JWT
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 返回用于验证签名的密钥
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	// 如果验证通过，返回解析后的 Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
