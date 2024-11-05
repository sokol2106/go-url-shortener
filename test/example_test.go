package test

import (
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/handlers"
	"github.com/sokol2106/go-url-shortener/internal/middleware"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func newServer() *httptest.Server {
	objectStorage := storage.NewMemory()
	srvShortURL := service.NewShortURL("http://localhost:8080", objectStorage)
	sh := handlers.NewHandlers(srvShortURL, middleware.NewToken(), "")
	server := httptest.NewServer(handlers.Router(sh))
	return server
}

func Example() {
	server := newServer()
	request, err := http.NewRequest(http.MethodPost, server.URL, strings.NewReader("https://ya.ru/"))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := server.Client().Do(request)
	if err != nil {
		return
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	defer response.Body.Close()
	fmt.Println(string(resBody))

	urlParse, err := url.Parse(string(resBody))
	if err != nil {
		return
	}

	request, err = http.NewRequest(http.MethodGet, server.URL, strings.NewReader(urlParse.Path))
	if err != nil {
		return
	}

	// PostJSON

	reqBody := "{\"url\": \"https://practicum.yandex.ru\"}"
	request, err = http.NewRequest(http.MethodPost, server.URL+"/api/shorten", strings.NewReader(reqBody))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err = server.Client().Do(request)
	if err != nil {
		return
	}
	resBody, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	defer response.Body.Close()
	fmt.Println(string(resBody))

	// PostBatch

	reqBody = "[{\"correlation_id\": \"1111\",\"original_url\": \"https://www.ozon.ru\"}," +
		"{\"correlation_id\": \"2222\",\"original_url\": \"https://ya.ru\"}]"
	request, err = http.NewRequest(http.MethodPost, server.URL+"/api/shorten/batch", strings.NewReader(reqBody))
	if err != nil {
		return
	}
	response, err = server.Client().Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	fmt.Println("ExamplePostBatch")
	fmt.Println(response.StatusCode)

	// GetPing

	request, err = http.NewRequest(http.MethodGet, server.URL+"/ping", strings.NewReader(""))
	if err != nil {
		return
	}
	response, err = server.Client().Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	fmt.Println("ExampleGetPing")
	fmt.Println(response.StatusCode)

	// GetUserShortenedURLs

	request, err = http.NewRequest(http.MethodGet, server.URL+"/api/user/urls", strings.NewReader(""))
	if err != nil {
		return
	}
	response, err = server.Client().Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	fmt.Println("ExampleGetUserShortenedURLs")
	fmt.Println(response.StatusCode)

	// DeleteUserShortenedURLs

	reqBody = "[\"rfgtyhju\",\"qazxswed\",\"wsxcderf\"]"
	request, err = http.NewRequest(http.MethodDelete, server.URL+"/api/user/urls", strings.NewReader(reqBody))
	if err != nil {
		return
	}

	response, err = server.Client().Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	fmt.Println("ExampleDeleteUserShortenedURLs")
	fmt.Println(response.StatusCode)
}
