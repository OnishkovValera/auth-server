package router

import (
	"auth-server/internal/http/handlers"
	mdw "auth-server/internal/http/middleware"
	"auth-server/internal/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(maker *jwt.Maker, handler *handlers.UserHandler) (*chi.Mux, error) {
	rt := chi.NewRouter()
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	rt.Use(middleware.StripSlashes)
	rt.Post("/register", handler.CreateUser)
	rt.Post("/login", handlers.JsonReturnJWT)
	rt.Group(func(r chi.Router) {
		r.Use(mdw.JWTAuthMiddleware(maker))
		r.Get("/info", handlers.MyInfo)
	})

	return rt, nil
}
