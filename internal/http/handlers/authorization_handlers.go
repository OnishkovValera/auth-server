package handlers

import (
	"auth-server/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AuthorizationHandler struct {
	s *service.AuthorizationService
}

func NewAuthorizationHandler(s *service.AuthorizationService) *AuthorizationHandler {
	return &AuthorizationHandler{s: s}
}

func (h *AuthorizationHandler) GetUserInfoHandler(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	rawToken, err := h.getRawToken(request)
	claims, err := h.s.GetClaimsFromToken(rawToken)
	if err != nil {
		fmt.Println(err)
	}

	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(claims)
	if err != nil {
		return
	}
}

func (h *AuthorizationHandler) getRawToken(r *http.Request) (string, error) {
	rawToken := r.Header.Get("Authorization")
	prefix := "Bearer "
	if strings.HasPrefix(rawToken, prefix) {
		rawToken = strings.TrimPrefix(rawToken, prefix)
		return rawToken, nil
	}
	NoBearerToken := errors.New("no bearer token found")
	return "", NoBearerToken

}
