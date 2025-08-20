package main

import (
	"todo/internal/log"
	"todo/internal/middleware"
	"todo/internal/http/handler"
	"todo/internal/storage/postgres"
	"todo/internal/storage/redis"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	db.StartDB()
	redis.StartRedis()

	router := handler.BaseHandler{}
	middleware := middleware.LoggingMiddleWare(middleware.AuthMiddleWare(router))

	log.Logger.Info("Listening on port " + port)
	http.ListenAndServe(port, middleware)
}
