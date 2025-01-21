package token

import (
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	jwt, err := GenerateJWT(1, "test", time.Minute)
	if err != nil {
		t.Errorf("GenerateJWT() error = %v", err)
	}
	t.Logf("GenerateJWT() = %v", jwt)
}

func TestParseJWT(t *testing.T) {
	jwt, err := GenerateJWT(1, "test", time.Minute)
	if err != nil {
		t.Errorf("GenerateJWT() error = %v", err)
	}
	claims, err := ParseJWT(jwt)
	if err != nil {
		t.Errorf("ParseJWT() error = %v", err)
	}
	t.Logf("ParseJWT() = %v", claims)
}
