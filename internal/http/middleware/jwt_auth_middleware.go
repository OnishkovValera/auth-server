package middleware

import (
	"auth-server/internal/service"
	"fmt"
	"net/http"

	"strings"
)

type JWTMiddleware struct {
	ah *service.AuthorizationService
}

func NewJWTMiddleware(ah *service.AuthorizationService) *JWTMiddleware {
	return &JWTMiddleware{ah: ah}
}

func (jw *JWTMiddleware) JWTAuthMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			prefix := "Bearer "
			if strings.HasPrefix(token, prefix) {
				token = strings.TrimPrefix(token, prefix)
			}

			checkJWT, err := jw.ah.CheckJWT(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Println(err)
					return
				}
				return
			}

			if !checkJWT {
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("Unauthorized"))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
