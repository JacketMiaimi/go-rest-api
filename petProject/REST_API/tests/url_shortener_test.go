package tests

import (
	"main.go/internal/http-server/handlers/url/save"
	"main.go/internal/lib/api"
	"main.go/internal/lib/random"
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6" // Генерирует случайные данные (email, pass)
	"github.com/gavv/httpexpect/v2"   // Нужен для тестирования http сервиса
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8080"
)

func TestURLShortener_HappyPath(t *testing.T) { // Для понимания метода
	u := url.URL{ // Базовый url который обращается клиент
		Scheme: "http",
		Host:   host, // 8080
	}
	e := httpexpect.Default(t, u.String()) // Создаем клиента в него будут добавляться запросы

	e.POST("/url"). // формируем url (потом будем работать с клиентом). post - формирует наш запрос
			WithJSON(save.Request{ // Берем объект json из нашего request
			URL:   gofakeit.URL(),             // генерация url
			Alias: random.NewRandomString(10), // Длина нашего сокр ссылки
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().            // что ожидать от ответа
		Status(200).         // ожидаем 200
		JSON().              // Формируем его в json
		Object().            // Из него получаем объект
		ContainsKey("alias") // Проверка, что наш объект содержит alias
	// Todo: можно добавить еще больше
}

//nolint:funlen
func TestURLShortener_SaveRedirect(t *testing.T) { // Сценарий тестирования. Сначала будет сохранять, редеректать и удалять
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "failed URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
		// TODO: add more test cases
	}

	for _, tc := range testCases { // проходится
		t.Run(tc.name, func(t *testing.T) { // Запускает наши кейсы
			u := url.URL{ // Базовый url который обращается клиент
				Scheme: "http",
				Host:   host, // 8080
			}

			e := httpexpect.Default(t, u.String()) // Создаем клиент с помощью его будут доб запросы

			// Save

			resp := e.POST("/url"). // Отправляем запрос на url | потом будем работать с тем клиентом | Post-формирует наш запрос
						WithJSON(save.Request{ // Передаем объект json из нашего request, который будет передаваться в json
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("myuser", "mypass").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias // alias тест кейсов

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias) // проверка, что в теле ответа именно alias, который мы отправили и он его вернет
			} else {
				resp.Value("alias").String().NotEmpty() // Здесь тоже возвращает

				alias = resp.Value("alias").String().Raw() // Сохраняем то что к нам пришло в alias
			}

			// Redirect

			testRedirect(t, alias, tc.url) // Проверяем редирект
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{ // формируем url
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String()) // делает HTTP-запрос и получает конечный URL (Location и 302 и тд)
	require.NoError(t, err)                             // проверка на ошибку

	require.Equal(t, urlToRedirect, redirectedToURL) // Проверяем, что полученный редирект совпадает с тем, что мы положили в аргументе
}
