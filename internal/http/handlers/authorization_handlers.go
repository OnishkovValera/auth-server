package handlers

import (
	"auth-server/internal/model/dto"
	"auth-server/internal/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthorizationHandler struct {
	validator *validator.Validate
	s         *service.AuthorizationService
}

func NewAuthorizationHandler(s *service.AuthorizationService, v *validator.Validate) *AuthorizationHandler {
	return &AuthorizationHandler{s: s, validator: v}
}

func (h *AuthorizationHandler) Login(writer http.ResponseWriter, request *http.Request) {
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

	token, err := h.s.Login(request.Context(), login)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	jsonToken, err := json.Marshal(map[string]string{"access_token": token.AccessToken})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     "refresh_token",
		HttpOnly: true,
		Path:     "/refresh",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Value:    token.RefreshToken,
		MaxAge:   60 * 60 * 24 * 90,
	})

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(jsonToken)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (h *AuthorizationHandler) Refresh(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("refresh_token")
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}
	token, err := h.s.RefreshToken(request.Context(), cookie.Value)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(err.Error()))
	}

	if token == nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("problem with token"))
		return
	}
	jsonToken, err := json.Marshal(map[string]string{"access_token": token.AccessToken})
	http.SetCookie(writer, &http.Cookie{
		Name:     "refresh_token",
		Path:     "/refresh",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Value:    token.RefreshToken,
		MaxAge:   60 * 60 * 24 * 90,
	})
	writer.Write(jsonToken)
	writer.WriteHeader(http.StatusOK)

}
