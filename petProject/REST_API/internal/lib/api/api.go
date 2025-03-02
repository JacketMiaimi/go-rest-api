// делает HTTP-запрос и получает конечный URL

package api

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidStatusCode = errors.New("invalid status code")
)

// GetRedirect returns the final URL after redirection.
func GetRedirect(url string) (string, error) { // 1 - принимает url по которому делает запрос | возвращает конечный url
	const op = "api.GetRedirect"

	client := &http.Client{ // Создается кастомный клиент
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // stop after 1st redirect
		}, // не следует за редиректами, просто return ответ
	}

	resp, err := client.Get(url) // Запрос отправляется
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s: %w: %d", op, ErrInvalidStatusCode, resp.StatusCode)
	}

	return resp.Header.Get("Location"), nil // возвращает HTTP-ответ с кодом редиректа (например, 302) и заголовком Location
}
