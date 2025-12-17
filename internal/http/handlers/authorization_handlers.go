package handlers

import (
	"auth-server/internal/jwt"
	"auth-server/internal/model/dto"
	"encoding/json"
	"fmt"
	"net/http"
)

func JsonReturnJWT(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := dto.NewUserDto()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		_ = fmt.Errorf("json decoding error: %s", err)
		return
	}

	maker := jwt.NewMaker()
	createJWT, err := maker.CreateJWT(user.Login)
	if err != nil {
		fmt.Println(err)
		return
	}
	jwtDto := dto.NewJWTDto(createJWT)
	err = json.NewEncoder(w).Encode(jwtDto)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CheckJWT(token string, maker *jwt.Maker) (bool, error) {
	if isValid, err := maker.ParseJWT(token); !isValid || err != nil {
		return false, err
	}
	return true, nil
}

func Register(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
}

func MyInfo(writer http.ResponseWriter, request *http.Request) {
}
