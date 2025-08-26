package log

import (
	"io"
	"log"
	"log/slog"
	"os"
	"todo/internal/config"
)

func NewLogger(cfg config.Config) (*slog.Logger, *os.File) {
	var writer io.Writer = os.Stdout
	var file *os.File

	if cfg.LogPath != "" {
		file, err := os.OpenFile(cfg.LogPath, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0644)

		if err != nil {
			log.Fatal("logger failed to initialize:", err)
		}
	
		writer = file
	}

	var level slog.Level

	switch cfg.LogLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})), file
}
