package handlers

import (
	"auth-server/internal/model/dto"
	"auth-server/internal/model/entity"
	"auth-server/internal/service"
	"encoding/json"
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

func (h *UserHandler) GetUserInfo(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	rawToken := request.Header.Get("Authorization")
	prefix := "Bearer "
	if strings.HasPrefix(rawToken, prefix) {
		rawToken = strings.TrimPrefix(rawToken, prefix)
	}
	claims, err := h.userService.GetDataFromUser(rawToken)
	if err != nil {
		fmt.Println(err)
	}

	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(claims)
	if err != nil {
		return
	}
}

func (h *UserHandler) Login(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	login := dto.NewLoginRequest()
	err := json.NewDecoder(request.Body).Decode(login)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	err = h.validator.Struct(login)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	token, err := h.userService.Login(request.Context(), login)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	jsonToken, err := json.Marshal(token)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(jsonToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (h *UserHandler) Refresh(writer http.ResponseWriter, request *http.Request) {

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
