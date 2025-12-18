package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func CreateClaims(login string) (*Claims, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Claims{
			Login: login,
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        token.String(),
				Issuer:    "onish-auth-service",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
		nil
}
