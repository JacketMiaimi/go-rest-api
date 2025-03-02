package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"main.go/internal/config"
	"main.go/internal/http-server/handlers/redirect"
	"main.go/internal/http-server/handlers/url/save"
	"main.go/internal/lib/logger/handlers/slogpretty"
	"main.go/internal/lib/logger/sl"
	"main.go/internal/storage/sqlite"
	"net/http"
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

	log.Info("starting url-shortener",
		slog.String("env", cfg.Env),
		//slog.String("version", "123"), // ?
	)
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// Инициализируем роутер. Устанавливаем пакет chi
	// middleware - это цепочки когда наш handler обрабатывает запрос и основной называется handler запроса, а другие middleware
	// Он проверяет авторизацию и не дает пройти, если неправильно

	router := chi.NewRouter()
	router.Use(middleware.RequestID) // Он добавляет к каждому запросу RequestID - это уникальный идентификатор для каждого запроса в системе (для ошибок, логирование )
	router.Use(middleware.Logger)    // логирует все входящие запросы
	//router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer) // Паника
	router.Use(middleware.URLFormat) // Красивые URL

	router.Route("/url", func(r chi.Router) { // 1 - общий префикс url
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		//r.Delete("/{alias}", Dellete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage)) // 1 - имя параметра, что бы дальше получить его в handler

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:              cfg.Address, // адрес из конфига
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPServer.Timeout, // Время на обработку запроса
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
	// TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger // Объявляем логер

	switch env {
	case envLocal:
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
