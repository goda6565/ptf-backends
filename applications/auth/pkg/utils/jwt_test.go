package utils

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupEnv(t *testing.T) {
	defer func() {
		err := os.Unsetenv("JWT_SECRET_KEY")
		assert.NoError(t, err, "Unsetenv should not return an error")
	}()
	t.Helper()
	err := os.Setenv("JWT_SECRET_KEY", "mysecret")
	assert.NoError(t, err, "Setting JWT_SECRET_KEY should not error")
}

func TestGenerateAndValidateToken(t *testing.T) {
	setupEnv(t)

	userID := uint(1)
	email := "Alice"

	// トークン生成
	tokenString, err := GenerateSignedString(userID, email)
	assert.NoError(t, err, "GenerateSignedString should not return error")
	assert.NotEmpty(t, tokenString, "Token string should not be empty")

	// トークン検証
	claims, err := ValidateToken(tokenString)
	assert.NoError(t, err, "ValidateToken should not return error")
	assert.Equal(t, int(userID), claims.ID, "User ID should match")
	assert.Equal(t, email, claims.Email, "email should match")
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	setupEnv(t)

	userID := uint(2)
	email := "Bob"

	// 正しいシークレットでトークン生成
	tokenString, err := GenerateSignedString(userID, email)
	assert.NoError(t, err, "GenerateSignedString should not return error")

	// シークレットキーを変更して、署名エラーを発生させる
	err = os.Setenv("JWT_SECRET_KEY", "wrongsecret")
	assert.NoError(t, err, "Setting wrong secret should not error")

	_, err = ValidateToken(tokenString)
	assert.Error(t, err, "ValidateToken should return error for invalid signature")
	assert.True(t, strings.Contains(err.Error(), "signature"), "Error message should contain 'signature'")
}

func TestValidateToken_Expired(t *testing.T) {
	setupEnv(t)

	// 有効期限が過去のトークンを手動で作成
	expiredTime := time.Now().Add(-time.Hour)
	claims := MyJWTClaims{
		ID:    "3",
		Email: "Charlie",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			Issuer:    "ptf-auth-service",
			Subject:   "3",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	assert.NoError(t, err, "Signing token should not error")

	_, err = ValidateToken(tokenString)
	assert.Error(t, err, "ValidateToken should return error for expired token")
	assert.True(t, strings.Contains(err.Error(), "expired"), "Error message should mention 'expired'") // エラーメッセージに 'expired' が含まれていること
}

func TestValidateToken_InvalidFormat(t *testing.T) {
	setupEnv(t)

	_, err := ValidateToken("not-a-valid-token")
	assert.Error(t, err, "ValidateToken should return error for invalid token format")
}
