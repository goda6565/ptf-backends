package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyJWTClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func GenerateSignedString(userId uint, email string) (string, error) {
	claims := MyJWTClaims{
		ID:    strconv.Itoa(int(userId)),
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ptf-auth-service",
			Subject:   strconv.Itoa(int(userId)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

type TokenClaims struct {
	ID    int
	Email string
}

func ValidateToken(signedToken string) (*TokenClaims, error) {
	// MyJWTClaims 型を使ってパース
	token, err := jwt.ParseWithClaims(signedToken, &MyJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムをチェック（HMAC 系のみ許可）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("signature validation failed")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token is expired")
		}
		return nil, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(*MyJWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("unauthorized")
	}

	// クレーム内の ID は文字列として格納しているため、整数に変換
	id, err := strconv.Atoi(claims.ID)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	return &TokenClaims{
		ID:    id,
		Email: claims.Email,
	}, nil
}
