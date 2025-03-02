package Dellete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	resp "main.go/internal/lib/api/response"
	"main.go/internal/lib/logger/sl"
	"main.go/internal/storage"
	"net/http"
)

type URLDelete interface {
	DelURL(alias string) (string, error)
}

func New(log *slog.Logger, urlDelete URLDelete) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Delete.delete.Del"

		log := log.With(
			slog.String("op", op),
			slog.String("delete_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid delete"))

			return
		}

		_, err := urlDelete.DelURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("Url delete", slog.String("url", alias))
		w.WriteHeader(http.StatusOK)
	}
}
