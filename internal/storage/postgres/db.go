package db

import (
	"database/sql"
	"todo/internal/log"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	DB          *sql.DB
	PG_HOST     = os.Getenv("PG_HOST")
	PG_USER     = os.Getenv("PG_USER")
	PG_PASSWORD = os.Getenv("PG_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
)

func StartDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		PG_HOST,
		PG_USER,
		PG_PASSWORD,
		DB_NAME,
	)

	var err error

	DB, err = sql.Open("postgres", dsn)

	if err != nil {
		log.Logger.Error("postgres connection failed", "err", err)
		os.Exit(1)
	}

	log.Logger.Info("postgres connection established")
}
