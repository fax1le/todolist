package log

import (
	"io"
	"log/slog"
	"os"
)

func NewLogger(w io.Writer, level slog.Leveler) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))
}

var Logger = NewLogger(os.Stdout, slog.LevelDebug)
