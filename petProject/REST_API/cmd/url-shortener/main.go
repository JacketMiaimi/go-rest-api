package main

import (
	_ "fmt"
	"log/slog"
	"main.go/internal/config"
	"main.go/internal/storage/lib/logger/sl"
	"main.go/internal/storage/sqlite"

	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	// Подключаем библиотеку cleanenv для валидации
	// Setenv - установить окружение
	os.Setenv("CONFIG_PATH", "C:\\IT\\backend\\Go\\petProject\\REST_API\\config\\local.yaml") //В переменную окружения CONFIG_PATH записывается путь до yaml

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	// TODO: init router: chi, "chi render"

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger // Объявляем логер

	switch env {
	case envLocal:
		log = slog.New(
			// 2 - это настройки логера{ будут выводится все логи, ведь минимальный уровень }
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}), // Обработчик (логирует)
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}), // Обработчик логов в формате json
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
