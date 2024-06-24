package test

import (
	"encoding/json"
	"fmt"
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

func TestHandlerMain(t *testing.T) {
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

func TestCheckPostHandlerMain(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		wantPost strWant
	}{
		{
			name: "Error httpss",
			url:  "localhost:8080",
			wantPost: strWant{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Error httpss",
			url:  "httpss://practicum.yandex.ru/",
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
		{
			name: "Error empty",
			url:  "",
			wantPost: strWant{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "success",
			url:  "https://yandex.ru/maps/15/tula/?ll=37.617348%2C54.193122&z=13",
			wantPost: strWant{
				code: http.StatusCreated,
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

func TestPostJSONHandlerMain(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantPost strWant
		wantGet  strWant
	}{
		{
			name: "Test POST JSON",
			body: "{\"url\": \"https://practicum.yandex.ru\"}",
			wantPost: strWant{
				code:        http.StatusCreated,
				contentType: "application/json",
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
			request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", server.URL, "/api/shorten"), strings.NewReader(tt.body))
			request.Header.Set("Content-Type", tt.wantPost.contentType)
			require.NoError(t, err)
			response, err := server.Client().Do(request)
			require.NoError(t, err)

			status := assert.Equal(t, tt.wantPost.code, response.StatusCode)
			content := assert.Equal(t, tt.wantPost.contentType, response.Header.Get("Content-Type"))
			if status && content {

				//
				resBody, err := io.ReadAll(response.Body)
				require.NoError(t, err)

				err = response.Body.Close()
				require.NoError(t, err)

				var respJS shorturl.ResponseJSON
				err = json.Unmarshal(resBody, &respJS)
				require.NoError(t, err)

				urlParse, err := url.Parse(respJS.Result)
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
