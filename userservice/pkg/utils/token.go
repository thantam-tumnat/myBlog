package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MyBlogsClaims is the payload carried by every access token.
// blogservice verifies tokens against the same structure + shared secret.
type MyBlogsClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken signs a token valid for 24h with the shared HS256 secret.
func GenerateAccessToken(secret string, userID uint, username, role string) (string, error) {
	claims := MyBlogsClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "userservice",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
