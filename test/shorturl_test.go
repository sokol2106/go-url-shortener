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
	objectStorage := storage.NewMemory()
	sh := shorturl.New("http://localhost:8080", objectStorage)
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

			// Попробую здесь testify/suite
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
	objectStorage := storage.NewMemory()
	sh := shorturl.New("http://localhost:8080", objectStorage)
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
	objectStorage := storage.NewMemory()
	sh := shorturl.New("http://localhost:8080", objectStorage)
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
	objectStorage := storage.NewFile(fileName)
	defer objectStorage.Close()

	sh := shorturl.New("http://localhost:8080", objectStorage)
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

				urlParse, err := url.Parse(string(resBody))
				require.NoError(t, err)

				resURL := objectStorage.GetURL(strings.ReplaceAll(urlParse.Path, "/", ""))
				assert.Equal(t, tt.url, resURL)

				err = sh.Close()
				require.NoError(t, err)

				server.Close()
				require.NoError(t, err)

				err = os.Remove(fileName)
				require.NoError(t, err)
			}
		})
	}
}

func TestShortURLTestify(t *testing.T) {
	objectStorage := storage.NewMemory()
	shrt := shorturl.New("http://localhost:8080", objectStorage)

	//  Post
	request := httptest.NewRequest("POST", "/", strings.NewReader("https://practicum.yandex.ru/"))
	response := httptest.NewRecorder()
	shrt.Post(response, request)
	assert.Equal(t, http.StatusCreated, response.Code)

	// PostJSON
	request = httptest.NewRequest("POST", "/", strings.NewReader("{\"url\": \"https://dzen.ru\"}"))
	request.Header.Set("Content-Type", "application/json")
	shrt.PostJSON(response, request)
	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))

	err := shrt.Close()
	require.NoError(t, err)

	// id ни как тут не получить
	// resBody, err := io.ReadAll(response.Body)
	// require.NoError(t, err)
	// request = httptest.NewRequest("GET", "/", strings.NewReader(string(resBody)))
	// shrt.Get(response, request)
	// assert.Equal(t, http.StatusOK, response.Code)

	// Попробую здесь testify/suite

}

func TestShortURLPostBatch(t *testing.T) {

	objectStorage := storage.NewMemory()
	shrt := shorturl.New("http://localhost:8080", objectStorage)

	t.Run("Test POST Batch", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/", strings.NewReader(""+
			"[{\"correlation_id\": \"1111\",\"original_url\": \"https://www.ozon.ru\"},"+
			"{\"correlation_id\": \"2222\",\"original_url\": \"https://ya.ru\"}]"))
		response := httptest.NewRecorder()
		shrt.PostBatch(response, request)
		assert.Equal(t, http.StatusCreated, response.Code)
		shrt.Close()

	})
}
