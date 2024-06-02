package shorturl

import (
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
				code:        http.StatusTemporaryRedirect,
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Проверяем Post запрос
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))

			w := httptest.NewRecorder()
			HanlerMain(w, request)
			res := w.Result()

			status := assert.Equal(t, tt.wantPost.code, res.StatusCode)

			content := assert.Equal(t, tt.wantPost.contentType, res.Header.Get("Content-Type"))
			if status && content {
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)

				urlParse, err := url.Parse(string(resBody))
				require.NoError(t, err)

				// Проверяем Get запрос
				request := httptest.NewRequest(http.MethodGet, urlParse.Path, strings.NewReader(tt.url))
				w := httptest.NewRecorder()
				HanlerMain(w, request)
				res := w.Result()

				assert.Equal(t, tt.wantGet.code, res.StatusCode)
				assert.Equal(t, tt.url, res.Header.Get("Location"))

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
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))

			w := httptest.NewRecorder()
			HanlerMain(w, request)
			res := w.Result()
			assert.Equal(t, tt.wantPost.code, res.StatusCode)
		})
	}
}
