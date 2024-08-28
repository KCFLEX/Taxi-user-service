package tokenservice

import (
	"context"
	"time"

	"github.com/KCFLEX/Taxi-user-service/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

type Token struct {
	SecretKey string
}

func New(config config.Config) *Token {
	return &Token{
		SecretKey: config.SecretKey,
	}
}

func (t *Token) GenerateToken(ctx context.Context, userID string, duration time.Duration) (string, error) {
	//token expiration time
	expirationTime := time.Now().Add(duration)

	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// create token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token
	tokenStr, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, err

}
