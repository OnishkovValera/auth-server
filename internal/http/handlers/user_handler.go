package handlers

import (
	"auth-server/internal/model/dto"
	"auth-server/internal/model/entity"
	"auth-server/internal/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	validator   *validator.Validate
	userService *service.UserService
}

func NewUserHandler(validate *validator.Validate, userService *service.UserService) *UserHandler {
	return &UserHandler{validator: validate, userService: userService}
}

func (h *UserHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	ctx := request.Context()
	userDto := dto.NewUserDto()
	if err := json.NewDecoder(request.Body).Decode(&userDto); err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
	}

	err := h.validator.Struct(userDto)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	passwordHash, err := service.HashPassword(userDto.Password)

	if err != nil {
		fmt.Println(err)
		return
	}

	user := entity.User{
		Login:        userDto.Login,
		Name:         userDto.Name,
		Surname:      userDto.Surname,
		PasswordHash: passwordHash,
		IsActive:     true,
	}

	err = h.userService.CreateUser(ctx, &user)
	if err != nil {
		return
	}

}
