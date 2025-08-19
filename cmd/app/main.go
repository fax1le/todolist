package main

import (
	"todo/internal/log"
	"todo/internal/storage/postgres"
	"todo/internal/storage/redis"
	"todo/internal/service/handler"
	"todo/internal/middleware"
	"net/http"
)

func main() {
	db.StartDB()
	redis.StartRedis()

	router := handler.BaseHandler{}
	middleware := middleware.LoggingMiddleWare(middleware.AuthMiddleWare(router))

	log.Logger.Info("Listening on port 8080...")
	http.ListenAndServe(":8080", middleware)
}
