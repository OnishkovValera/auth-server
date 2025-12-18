package router

import (
	"auth-server/internal/http/handlers"
	mdw "auth-server/internal/http/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type MainRouter struct {
	authorizationHandler *handlers.AuthorizationHandler
	userHandler          *handlers.UserHandler
	jwtMdw               *mdw.JWTMiddleware
}

func NewRouter(authHandler *handlers.AuthorizationHandler, userHandler *handlers.UserHandler, jwtMdw *mdw.JWTMiddleware) *MainRouter {
	return &MainRouter{authorizationHandler: authHandler, userHandler: userHandler, jwtMdw: jwtMdw}
}

func (mRouter *MainRouter) GetMux() *chi.Mux {
	rt := chi.NewRouter()
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	rt.Use(middleware.StripSlashes)
	rt.Post("/register", mRouter.userHandler.CreateUser)
	rt.Post("/login", mRouter.userHandler.Login)
	rt.Post("/refresh", mRouter.userHandler.Refresh)
	rt.Group(func(r chi.Router) {
		r.Use(mRouter.jwtMdw.JWTAuthMiddleware())
		r.Get("/info", mRouter.authorizationHandler.GetUserInfoHandler)
	})

	return rt
}
