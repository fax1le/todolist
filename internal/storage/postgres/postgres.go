package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"todo/internal/config"

	_ "github.com/lib/pq"
)

func StartDB(cfg config.Config, logger *slog.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PGHost,
		cfg.PGUser,
		cfg.PGPassword,
		cfg.DBName,
	)

	DB, err := sql.Open("postgres", dsn)

	if err != nil {
		return DB, err
	}

	err = DB.Ping()

	return DB, err
}
