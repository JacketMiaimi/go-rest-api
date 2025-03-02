package redirect_test

import (
	"main.go/internal/http-server/handlers/redirect"
	"main.go/internal/http-server/handlers/redirect/mocks"
	"main.go/internal/lib/api"
	"main.go/internal/lib/logger/handlers/slogdiscard"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success", // Здесь мы пытаемся поломать его
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
	}

	for _, tc := range cases { // проходимся по кейсу
		t.Run(tc.name, func(t *testing.T) { // t.run - запускает код с названием tc.name
			urlGetterMock := mocks.NewURLGetter(t) // Создаем объект мока

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias). // возвращает ошибку
					Return(tc.url, tc.mockError).Once() // Once - 1 раз
			}

			r := chi.NewRouter() // Создает роутер - это компонент, который распределяет HTTP-запросы на обработчики.
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))
			//Создаёт роут (/{alias}) для GET-запросов.
			//Назначает обработчик redirect.New(...), который выполняется, когда приходит запрос

			ts := httptest.NewServer(r) // Локальный http-сервис для тестирования
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias) // отправляет get запрос
			require.NoError(t, err)                                          // проверка на ошибку

			// Check the final URL after redirection.
			assert.Equal(t, tc.url, redirectedToURL) // сравнение tc.url и redirectedToURL
		})
	}
}
