package main

import (
	"github.com/avealice/filmhub/internal/app"
	_ "github.com/avealice/filmhub/internal/handler"
	_ "github.com/avealice/filmhub/internal/model"

	_ "github.com/avealice/filmhub/docs"
)

// @title FilmHub API
// @description API для работы с фильмами и актерами в фильмотеке FilmHub. Логин и пароль от админки: admin и kek.

// @contact.email avanesova_alisa@mail.ru
// @contact.url http://www.github.com/avealice
// @contact.telegram @ohhalice

// @host 127.0.0.1:8000
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	a := app.NewApp()
	a.Run()
}
