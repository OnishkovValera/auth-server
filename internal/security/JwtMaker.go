package security

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtMaker struct {
	secretKey string
}

func NewMaker(key string) *JwtMaker {
	return &JwtMaker{key}
}

func (jwtMaker *JwtMaker) CreateAccessToken(login string) (string, error) {
	claims, err := CreateClaims(login)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtMaker.secretKey))
	if err != nil {
		panic(err)
	}
	return tokenString, nil
}

func (jwtMaker *JwtMaker) CreateRefreshToken() string {
	return uuid.New().String()
}

func (jwtMaker *JwtMaker) ParseJWT(tokenString string) (bool, error) {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtMaker.secretKey), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (jwtMaker *JwtMaker) GetClaimsFromToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtMaker.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, err
}
