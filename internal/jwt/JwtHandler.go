package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var key = "89shf84928bv39"

type Maker struct {
	secretKey string
}

func NewMaker() *Maker {
	return &Maker{key}
}

func (jwtMaker *Maker) CreateJWT(login string) (string, error) {
	claims, err := CreateClaims("test@login.pidor")
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

func (jwtMaker *Maker) ParseJWT(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtMaker.secretKey), nil
	})
	fmt.Println(token, err)
	if err != nil {
		return false, err
	}
	return true, nil
}
