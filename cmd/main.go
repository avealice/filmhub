package main

import (
	"filmhub/internal/app"
	_ "filmhub/internal/handler"
	_ "filmhub/internal/model"

	_ "filmhub/docs"
)

// @title FilmHub API
// @description API для работы с фильмами и актерами в FilmHub. login: admin, password: kek

// @contact.url http://www.github.com/avealice

// @host 127.0.0.1:8000
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
// @BearerFormat Bearer
func main() {
	a := app.NewApp()
	a.Run()
}
