package test

import (
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type strWant struct {
	code        int
	body        string
	contentType string
}

func TestHanlerMain(t *testing.T) {

	tests := []struct {
		name     string
		url      string
		wantPost strWant
		wantGet  strWant
	}{
		{
			name: "Test redirect",
			url:  "https://practicum.yandex.ru/",
			wantPost: strWant{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
			wantGet: strWant{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Проверяем Post запрос

			server := httptest.NewServer(shorturl.ShortRouter("http://localhost:8080"))
			defer server.Close()
			request, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(tt.url))
			require.NoError(t, err)

			response, err := server.Client().Do(request)
			require.NoError(t, err)

			status := assert.Equal(t, tt.wantPost.code, response.StatusCode)
			content := assert.Equal(t, tt.wantPost.contentType, response.Header.Get("Content-Type"))

			if status && content {
				resBody, err := io.ReadAll(response.Body)
				require.NoError(t, err)

				err = response.Body.Close()
				require.NoError(t, err)

				urlParse, err := url.Parse(string(resBody))
				require.NoError(t, err)

				// Проверяем Get запрос

				request, err = http.NewRequest(http.MethodGet, server.URL+urlParse.Path, nil)
				require.NoError(t, err)

				response, err = server.Client().Do(request)
				require.NoError(t, err)

				assert.Equal(t, tt.wantGet.code, response.StatusCode)

				err = response.Body.Close()
				require.NoError(t, err)

			}
		})
	}
}

func TestErrorPostHanlerMain(t *testing.T) {

	tests := []struct {
		name     string
		url      string
		wantPost strWant
	}{
		{
			name: "Error httpss",
			url:  "httpss://practicum.yandex.ru/",
			wantPost: strWant{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Error not http or https",
			url:  "practicum.yandex.ru",
			wantPost: strWant{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Error not host",
			url:  "http://",
			wantPost: strWant{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Проверяем Post запрос

			server := httptest.NewServer(shorturl.ShortRouter("http://localhost:8080"))
			defer server.Close()
			request, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(tt.url))
			require.NoError(t, err)

			response, err := server.Client().Do(request)
			require.NoError(t, err)
			assert.Equal(t, tt.wantPost.code, response.StatusCode)

			err = response.Body.Close()
			require.NoError(t, err)
		})
	}
}
