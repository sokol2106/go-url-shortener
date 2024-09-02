package test

import (
	"context"
	"github.com/sokol2106/go-url-shortener/internal/handlers"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

type strWant struct {
	code        int
	body        string
	contentType string
}

func TestFileReadWrite(t *testing.T) {
	fileName := "hort-url-db.json"
	objectStorage := storage.NewFile(fileName)
	defer objectStorage.Close()

	srvShortURL := service.NewShortURL("http://localhost:8080", objectStorage)
	sh := handlers.NewHandlers(srvShortURL)
	server := httptest.NewServer(handlers.Router(sh))

	tests := []struct {
		name     string
		url      string
		wantPost strWant
	}{
		{
			name: "Test file Read/Write",
			url:  "https://practicum.yandex.ru/",
			wantPost: strWant{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Проверяем Post запрос
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

				resURL, err := objectStorage.GetOriginalURL(context.Background(), strings.ReplaceAll(urlParse.Path, "/", ""))
				require.NoError(t, err)
				assert.Equal(t, tt.url, resURL.OriginalURL)

				err = srvShortURL.Close()
				require.NoError(t, err)

				server.Close()
				require.NoError(t, err)

				err = os.Remove(fileName)
				require.NoError(t, err)
			}
		})
	}
}

func TestShortURLPostBatch(t *testing.T) {
	objectStorage := storage.NewMemory()
	srvShortURL := service.NewShortURL("http://localhost:8080", objectStorage)
	handler := handlers.NewHandlers(srvShortURL)

	t.Run("Test POST Batch", func(t *testing.T) {
		t.Parallel()
		request := httptest.NewRequest("POST", "/", strings.NewReader(""+
			"[{\"correlation_id\": \"1111\",\"original_url\": \"https://www.ozon.ru\"},"+
			"{\"correlation_id\": \"2222\",\"original_url\": \"https://ya.ru\"}]"))
		response := httptest.NewRecorder()
		handler.PostBatch(response, request)
		assert.Equal(t, http.StatusCreated, response.Code)
		srvShortURL.Close()

	})
}

func TestGetUserShortenedURLs(t *testing.T) {
	objectStorage := storage.NewMemory()
	srvShortURL := service.NewShortURL("http://localhost:8080", objectStorage)
	sh := handlers.NewHandlers(srvShortURL)
	server := httptest.NewServer(handlers.Router(sh))

	defer server.Close()

	t.Run("testGetUserShortenedURLs", func(t *testing.T) {
		// No Cooke
		request, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader("https://www.ozon.ru/"))
		require.NoError(t, err)

		response, err := server.Client().Do(request)
		require.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))
		require.Len(t, response.Cookies(), 1)
		cookieUser1 := response.Cookies()[0]

		// No Cooke
		request, err = http.NewRequest(http.MethodPost, server.URL, strings.NewReader("https://calendar.google.com/"))
		require.NoError(t, err)

		response, err = server.Client().Do(request)
		require.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))
		require.Len(t, response.Cookies(), 1)

		// Cooke
		request, err = http.NewRequest(http.MethodPost, server.URL, strings.NewReader("https://ya.ru/"))
		require.NoError(t, err)
		request.AddCookie(response.Cookies()[0])

		response, err = server.Client().Do(request)
		require.NoError(t, err)
		defer response.Body.Close()

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, "text/plain", response.Header.Get("Content-Type"))
		require.Len(t, response.Cookies(), 1)

		// user 2
		request, err = http.NewRequest(http.MethodGet, server.URL+"/api/user/urls", strings.NewReader(""))
		require.NoError(t, err)
		request.AddCookie(response.Cookies()[0])

		response, err = server.Client().Do(request)
		require.NoError(t, err)
		defer response.Body.Close()

		//_, err = io.ReadAll(response.Body)
		//require.NoError(t, err)
		//assert.Equal(t, "", bodyBytes)

		// user 1
		request, err = http.NewRequest(http.MethodGet, server.URL+"/api/user/urls", strings.NewReader(""))
		require.NoError(t, err)
		request.AddCookie(cookieUser1)

		response, err = server.Client().Do(request)
		require.NoError(t, err)
		defer response.Body.Close()

		//_, err = io.ReadAll(response.Body)
		//require.NoError(t, err)
		//assert.Equal(t, "", bodyBytes)

		// user 777
		request, err = http.NewRequest(http.MethodGet, server.URL+"/api/user/urls", strings.NewReader(""))
		require.NoError(t, err)
		auth := service.NewAuthorization()
		tkn, err := auth.NewUserToken()
		require.NoError(t, err)
		newCookie := http.Cookie{Name: "user", Value: tkn}
		request.AddCookie(&newCookie)

		response, err = server.Client().Do(request)
		require.NoError(t, err)
		defer response.Body.Close()

		//_, err = io.ReadAll(response.Body)
		//require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)

	})
}
