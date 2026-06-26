package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// MyBlogsClaims mirrors the token issued by userservice.
type MyBlogsClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// ParseAccessToken verifies the HS256 signature with the shared secret and
// returns the claims. It rejects tokens signed with an unexpected algorithm.
func ParseAccessToken(secret, tokenString string) (*MyBlogsClaims, error) {
	claims := &MyBlogsClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
