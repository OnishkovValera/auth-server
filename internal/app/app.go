package app

import (
	"auth-server/internal/config"
	"auth-server/internal/db"
	"auth-server/internal/http/handlers"
	mdw "auth-server/internal/http/middleware"
	"auth-server/internal/http/router"
	"auth-server/internal/repository"
	"auth-server/internal/security"
	"auth-server/internal/service"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type App struct {
	cnf config.Config
}

func NewApp(configPath string) *App {
	confPtr, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	return &App{cnf: *confPtr}
}

func (app *App) Run() {
	jwtMaker := security.NewMaker(app.cnf.JwtKey)
	vld := validator.New()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		app.cnf.DBUser, app.cnf.DBPassword, app.cnf.DBHost, app.cnf.DBPort, app.cnf.DBName)

	postgresDB, err := db.CreateDB(dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	tokenRepo := repository.NewRefreshTokenRepository(postgresDB)
	userRepo := repository.NewPostgresUserRepository(postgresDB)
	authService := service.NewAuthorizationService(jwtMaker, tokenRepo, userRepo)
	userService := service.NewUserService(userRepo, authService)
	userHandler := handlers.NewUserHandler(vld, userService)
	authHandler := handlers.NewAuthorizationHandler(authService)
	jwtMiddleware := mdw.NewJWTMiddleware(authService)
	rt := router.NewRouter(authHandler, userHandler, jwtMiddleware)

	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf(":%s", app.cnf.Port))
	err = http.ListenAndServe(fmt.Sprintf(":%s", app.cnf.Port), rt.GetMux())
	if err != nil {
		fmt.Println(err)
		panic(err)
		return
	}
}
