package app

import (
	"auth-server/internal/config"
	"auth-server/internal/db"
	"auth-server/internal/http/handlers"
	"auth-server/internal/http/router"
	"auth-server/internal/jwt"
	"auth-server/internal/repository"
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
	jwtMaker := jwt.NewMaker()
	vld := validator.New()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		app.cnf.DBUser, app.cnf.DBPassword, app.cnf.DBHost, app.cnf.DBPort, app.cnf.DBName)

	postgresDB, err := db.CreateDB(dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	userRepo := repository.NewPostgresUserRepository(postgresDB)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(vld, userService)
	rt, err := router.NewRouter(jwtMaker, userHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf(":%s", app.cnf.Port))
	err = http.ListenAndServe(fmt.Sprintf(":%s", app.cnf.Port), rt)
	if err != nil {
		fmt.Println(err)
		panic(err)
		return
	}
}
