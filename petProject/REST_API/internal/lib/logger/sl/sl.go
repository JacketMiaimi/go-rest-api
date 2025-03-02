package sl

import (
	_ "github.com/mattn/go-sqlite3" // Драйвер для sqlite
	"log/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

// Выводит ошибку
