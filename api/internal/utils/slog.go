package utils

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Slogger() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	Log = slog.New(handler)
}
