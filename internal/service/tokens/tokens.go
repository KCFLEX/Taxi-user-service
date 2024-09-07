package tokens

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
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

func (t *Token) ValidateToken(ctx context.Context, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(tt *jwt.Token) (interface{}, error) {
		if _, ok := tt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tt.Header["alg"])
		}
		return []byte(t.SecretKey), nil
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	if !token.Valid {
		return errorpac.ErrInvaiidToken
	}

	return nil
}

func (t *Token) ParseToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Subject, nil

}
