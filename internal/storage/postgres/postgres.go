package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"todo/internal/config"

	_ "github.com/lib/pq"
)

func StartDB(cfg config.Config, logger *slog.Logger) *sql.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PGHost,
		cfg.PGUser,
		cfg.PGPassword,
		cfg.DBName,
	)

	DB, err := sql.Open("postgres", dsn)

	if err != nil {
		logger.Error("postgres connection failed", "err", err)
		os.Exit(1)
	}

	logger.Info("postgres connection established")
	return DB
}
