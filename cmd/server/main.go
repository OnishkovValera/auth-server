package main

import (
	"auth-server/internal/app"
)

func main() {
	server := app.NewApp("internal/config/config.yaml")
	server.Run()
}
