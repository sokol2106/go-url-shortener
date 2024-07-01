package test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
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

func TestShortURL(t *testing.T) {
	var strg storage.ShortDataList
	strg.Init("")
	sh := shorturl.New("http://localhost:8080", strg)
	server := httptest.NewServer(shorturl.Router(sh))
	defer server.Close()

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

			// Не получается использовать httptest.NewRequest т.к. в shorturl не могу получить id
			// request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))
			// w := httptest.NewRecorder()
			// sh := shorturl.NewShortURL("http://localhost:8080")
			// sh.Post(w, request)
			// response := w.Result()

			// Как я понял httptest.NewServer не взаимодействует с httptest.NewRequest,
			// выпадает ошибка http: Request.RequestURI can't be set in client requests.

			// Сейчас у меня стартует сервер httptest.NewServer(shorturl.ShortRouter("http://localhost:8080"))
			// и я его тестирую запросами http.NewRequest

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

func TestPostJSON(t *testing.T) {
	var strg storage.ShortDataList
	strg.Init("")
	sh := shorturl.New("http://localhost:8080", strg)
	server := httptest.NewServer(shorturl.Router(sh))
	defer server.Close()

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
			request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", server.URL, "/api/shorten"), strings.NewReader(tt.body))
			request.Header.Set("Content-Type", tt.wantPost.contentType)
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

func TestGzipCompression(t *testing.T) {
	var strg storage.ShortDataList
	strg.Init("")
	sh := shorturl.New("http://localhost:8080", strg)
	server := httptest.NewServer(shorturl.Router(sh))
	defer server.Close()

	tests := struct {
		name     string
		url      string
		wantPost strWant
		wantGet  strWant
	}{
		name: "Test redirect GZIP",
		url:  "https://practicum.yandex.ru/",
		wantPost: strWant{
			code:        http.StatusCreated,
			contentType: "text/plain",
		},
		wantGet: strWant{
			code:        http.StatusOK,
			contentType: "text/plain",
		},
	}

	t.Run(tests.name, func(t *testing.T) {
		encodeBuffer := bytes.NewBuffer(nil)
		writeGZip := gzip.NewWriter(encodeBuffer)

		_, err := writeGZip.Write([]byte(tests.url))
		require.NoError(t, err)

		err = writeGZip.Close()
		require.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, server.URL, encodeBuffer)
		request.Header.Set("Content-Encoding", "gzip")
		request.Header.Set("Accept-Encoding", "gzip")
		require.NoError(t, err)

		response, err := server.Client().Do(request)
		require.NoError(t, err)

		status := assert.Equal(t, tests.wantPost.code, response.StatusCode)
		content := assert.Equal(t, tests.wantPost.contentType, response.Header.Get("Content-Type"))

		if status && content {
			readGZip, err := gzip.NewReader(response.Body)
			require.NoError(t, err)

			err = response.Body.Close()
			require.NoError(t, err)

			bodyDecode, err := io.ReadAll(readGZip)
			require.NoError(t, err)

			urlParse, err := url.Parse(string(bodyDecode))
			require.NoError(t, err)

			// Проверяем Get запрос

			request, err = http.NewRequest(http.MethodGet, server.URL+urlParse.Path, nil)
			require.NoError(t, err)
			response, err = server.Client().Do(request)
			require.NoError(t, err)

			assert.Equal(t, tests.wantGet.code, response.StatusCode)
			err = response.Body.Close()
			require.NoError(t, err)

		}
	})
}

func TestFileReadWrite(t *testing.T) {
	fileName := "hort-url-db.json"
	var strg storage.ShortDataList
	strg.Init(fileName)
	sh := shorturl.New("http://localhost:8080", strg)
	server := httptest.NewServer(shorturl.Router(sh))

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

				server.Close()
				err = sh.Close()
				require.NoError(t, err)

				urlParse, err := url.Parse(string(resBody))
				require.NoError(t, err)

				str := storage.ShortDataList{}
				str.Init(fileName)

				resURL := str.GetURL(strings.ReplaceAll(urlParse.Path, "/", ""))
				assert.Equal(t, tt.url, resURL)

				err = str.Close()
				require.NoError(t, err)

				err = os.Remove(fileName)
				require.NoError(t, err)
			}
		})
	}
}
