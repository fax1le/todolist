package main

import (
	"todo/internal/app"
	"todo/internal/config"
)

func main() {
	cfg := config.Load()
	app := application.New(cfg)
	app.Init()
	app.Run()
}
