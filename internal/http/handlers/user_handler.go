package handlers

import (
	"auth-server/internal/model/dto"
	"auth-server/internal/model/entity"
	"auth-server/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	var userDto *dto.UserDto
	var err error

	if userDto, err = h.decodeRequest(writer, request); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	passwordHash, err := service.HashPassword(userDto.Password)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
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
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(strconv.Itoa(user.Id)))
}

func (h *UserHandler) GetUserInfoHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	rawToken, err := h.getRawToken(request)
	user, err := h.userService.GetUserInfo(request.Context(), rawToken)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	userDto := dto.NewUserDtoWithParams(user.Login, user.Name, user.Surname)

	body, err := json.Marshal(userDto)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(body)
}

func (h *UserHandler) getRawToken(r *http.Request) (string, error) {
	rawToken := r.Header.Get("Authorization")
	prefix := "Bearer "
	if strings.HasPrefix(rawToken, prefix) {
		rawToken = strings.TrimPrefix(rawToken, prefix)
		return rawToken, nil
	}
	NoBearerToken := errors.New("no bearer token found")
	return "", NoBearerToken

}

func (h *UserHandler) decodeRequest(writer http.ResponseWriter, request *http.Request) (*dto.UserDto, error) {
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
		return nil, err
	}
	return userDto, nil
}
