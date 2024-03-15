package main

import (
	"filmhub/internal/app"
	_ "filmhub/internal/handler"
	_ "filmhub/internal/model"

	_ "filmhub/docs"
)

// @title FilmHub API
// @description API для работы с фильмами и актерами в FilmHub.
// @termsOfService http://swagger.io/terms/

// @contact.url http://www.github.com/avealice

// @host localhost:8000
func main() {
	a := app.NewApp()
	a.Run()
}
