package save

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	resp "main.go/internal/lib/api/response"
	"main.go/internal/lib/logger/sl"
	"main.go/internal/lib/random"
	"main.go/internal/storage"
	"net/http"
)

// Эта структура для парсинга выходящих данных json

type Request struct { // Структура запроса для парсинга
	URL   string `json:"url" validate:"url"` // обязательное поле для валидного URL
	Alias string `json:"alias,omitempty"`    // omitempty - если nil. То есть мы указываем, что должно отобразится
	//необязательное поле для псевдонима. Если оно не указано, оставляется пустым.
}

// Это структура для формирования ответа с дополнительным полем Alias.

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"` // Alias, который был использован
}

// TODO: move to config (перенести в конфиг)
const aliasLength = 6

//var w http.ResponseWriter

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	// который будет сохранять URL с псевдонимом в базе данных.
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
// New возвращает функцию-обработчик HTTP-запроса, которая будет вызываться при обработке запроса POST

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		// Чтение тела запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to read request body: "+err.Error()))
			return
		}
		defer r.Body.Close()

		// Если тело пустое, создаем пустую структуру, но не декодируем его
		if len(body) == 0 {
			log.Info("Request body is empty, proceeding with default values")
		} else {
			// Парсим JSON в структуру
			err = json.Unmarshal(body, &req)
			if err != nil {
				log.Error("failed to decode request body", sl.Err(err))
				render.JSON(w, r, resp.Error("failed to decode request: "+err.Error()))
				return
			}
			log.Info("Decoded JSON: ", req)
		}

		// Если URL пустой, генерируем его
		if req.URL == "" {
			req.URL = "https://generated-url.com/" + random.NewRandomString(10)
			log.Info("URL was empty, generated new one: ", slog.String("url", req.URL))
		}

		// Валидация
		if err := validator.New().Struct(req); err != nil {
			// Преобразуем ошибку валидации в тип validator.ValidationErrors
			validateErrs := err.(validator.ValidationErrors)

			// Логируем ошибку
			log.Error("invalid request", sl.Err(err))

			// Создаем ответ с ошибками валидации
			render.JSON(w, r, resp.ValidationError(validateErrs)) // Используем ValidationError с подробной информацией
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(6)
			log.Info("Generated alias: ", slog.String("alias", alias))
		}

		// Обработка сохранения URL
		id, err := urlSaver.SaveURL(req.URL, alias)
		// SaveURL - сохр url и alias в бд и возвр ID, ошибку
		// 1 - url наш | 2 - имя (сокращенное имя url)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("URL already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists")) // Создание ответа Json клиенту
			return
		}
		if err != nil {
			log.Error("failed to add URL", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add URL"))
			return
		}

		log.Info("URL added", slog.Int64("id", id))
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
