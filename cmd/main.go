package main

import (
	"filmhub/internal/app"

	_ "filmhub/docs"
)

func main() {

	a := app.NewApp()
	a.Run()
}
