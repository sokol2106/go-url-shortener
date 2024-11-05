// Package handlers предоставляет обработчики для HTTP-запросов, связанных с сокращением URL.
// Он включает в себя методы для создания, получения и удаления сокращенных URL,
// а также обработку запросов в формате JSON и пакетной обработки.
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sokol2106/go-url-shortener/internal/cerrors"
	"github.com/sokol2106/go-url-shortener/internal/middleware"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"io"
	"log"
	"net"
	"net/http"
)

// Handlers представляет собой структуру, содержащую сервисы для обработки URL и авторизации.
type Handlers struct {
	srvShortURL      *service.ShortURL      // Сервис сокращения URL
	token            *middleware.Token      // Токен авторизации
	srvAuthorization *service.Authorization // Сервис авторизации
	trustedSubnet    string                 // Конфигурация сервера
}

// RequestJSON представляет структуру запроса в формате JSON для создания сокращенного URL.
type RequestJSON struct {
	URL string `json:"url"`
}

// ResponseJSON представляет структуру ответа в формате JSON с сокращенным URL.
type ResponseJSON struct {
	Result string `json:"result"`
}

// ResponseStats представляет структуру ответа GetStats
type ResponseStats struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}

// NewHandlers создает новый экземпляр Handlers с переданным сервисом сокращения URL.
func NewHandlers(srv *service.ShortURL, t *middleware.Token, subnet string) *Handlers {
	srvAu := t.GetAuthorization()
	return &Handlers{
		srvShortURL:      srv,
		token:            t,
		srvAuthorization: srvAu,
		trustedSubnet:    subnet,
	}
}

// handlerError обрабатывает ошибки и возвращает соответствующий код состояния HTTP.
// Следующие коды могут вернуться:
// - 400 Bad Request: для всех прочих ошибок.
// - 409 Conflict: если пытаетесь добавить уже существующий оригинальный URL.
// - 410 Gone: если URL был помечен как удаленный.
func (h *Handlers) handlerError(err error) int {
	statusCode := http.StatusBadRequest
	if errors.Is(err, cerrors.ErrNewShortURL) {
		statusCode = http.StatusConflict
	}
	if errors.Is(err, cerrors.ErrGetShortURLDelete) {
		statusCode = http.StatusGone
	}

	log.Printf("error handling request: %v, status: %d", err, statusCode)
	return statusCode
}

// Post обрабатывает POST-запрос для создания сокращенного URL.
func (h *Handlers) Post(w http.ResponseWriter, r *http.Request) {
	handlerStatus := http.StatusCreated
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		handlerStatus = h.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	shortURL, err := h.srvShortURL.AddOriginalURL(string(body), h.srvAuthorization.GetCurrentUserID())

	if err != nil {
		handlerStatus = h.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(handlerStatus)
	w.Write([]byte(shortURL))
}

// PostJSON обрабатывает POST-запрос для создания сокращенного URL в формате JSON.
func (h *Handlers) PostJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handlerStatus := http.StatusCreated
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		handlerStatus = h.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	var (
		reqJS RequestJSON
		resJS ResponseJSON
	)

	err = json.Unmarshal(body, &reqJS)
	if err != nil {
		handlerStatus = h.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	resJS.Result, err = h.srvShortURL.AddOriginalURL(reqJS.URL, h.srvAuthorization.GetCurrentUserID())

	if err != nil {
		handlerStatus = h.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	resBody, err := json.Marshal(resJS)
	if err != nil {
		handlerStatus = h.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(handlerStatus)
	w.Write(resBody)
}

// PostBatch обрабатывает POST-запрос для пакетного создания сокращенных URL.
func (h *Handlers) PostBatch(w http.ResponseWriter, r *http.Request) {
	var (
		requestBatch  []service.RequestBatch
		responseBatch []service.ResponseBatch
		resBody       bytes.Buffer
		err           error
	)

	handlerStatus := http.StatusCreated

	if err = json.NewDecoder(r.Body).Decode(&requestBatch); err != nil {
		handlerStatus = h.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	responseBatch, err = h.srvShortURL.AddOriginalURLBatch(requestBatch, h.srvAuthorization.GetCurrentUserID())

	if err != nil {
		handlerStatus = h.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	err = json.NewEncoder(&resBody).Encode(responseBatch)
	if err != nil {
		handlerStatus = h.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	bodyB, err := io.ReadAll(&resBody)
	if err != nil {
		handlerStatus = h.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(handlerStatus)
	w.Write(bodyB)
}

// Get обрабатывает GET-запрос для получения оригинального URL по сокращенному ID.
func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	path := chi.URLParam(r, "id")
	URL, err := h.srvShortURL.GetOriginalURL(ctx, path)

	if err != nil {
		w.WriteHeader(h.handlerError(err))
		return
	}

	w.Header().Set("Location", URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// GetAll обрабатывает GET-запрос для получения всех сокращенных URL.
// По умолчанию возвращает статус 400.
func (h *Handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

// GetPing обрабатывает GET-запрос для проверки работоспособности сервиса.
func (h *Handlers) GetPing(w http.ResponseWriter, r *http.Request) {
	err := h.srvShortURL.PingContext()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetUserShortenedURLs обрабатывает GET-запрос для получения сокращенных URL текущего пользователя.
func (h *Handlers) GetUserShortenedURLs(w http.ResponseWriter, r *http.Request) {
	if h.srvAuthorization.IsNewUser() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	res, err := h.srvShortURL.GetUserShortenedURLs(ctx, h.srvAuthorization.GetCurrentUserID())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(res) < 4 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// DeleteUserShortenedURLs обрабатывает DELETE-запрос для удаления сокращенных URL текущего пользователя.
func (h *Handlers) DeleteUserShortenedURLs(w http.ResponseWriter, r *http.Request) {
	if h.srvAuthorization.IsNewUser() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var (
		request []string
	)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(h.handlerError(err))
		return

	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	h.srvShortURL.DeleteOriginalURLs(ctx, h.srvAuthorization.GetCurrentUserID(), request)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

// GetStats обрабатывает GET-запрос для получения количество сокращённых URL и пользователей в сервисе
func (h *Handlers) GetStats(w http.ResponseWriter, r *http.Request) {
	if h.trustedSubnet == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ipStr := r.Header.Get("X-Real-IP")
	ip := net.ParseIP(ipStr)
	if ip == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	_, cidr, err := net.ParseCIDR(h.trustedSubnet)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !cidr.Contains(ip) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var res ResponseStats

	res.Users = h.srvAuthorization.GetUsers()
	res.Urls = h.srvShortURL.GetURLs()

	resBody, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(h.handlerError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)

}

// Router создает маршрутизатор с заданными обработчиками и middleware.
func Router(handler *Handlers) chi.Router {
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.СompressionResponseRequest)
	router.Use(middleware.LoggingResponseRequest)
	router.Use(handler.token.TokenResponseRequest)

	// router
	router.Post("/", http.HandlerFunc(handler.Post))
	router.Post("/api/shorten", http.HandlerFunc(handler.PostJSON))
	router.Post("/api/shorten/batch", http.HandlerFunc(handler.PostBatch))
	router.Get("/{id}", http.HandlerFunc(handler.Get))
	router.Get("/*", http.HandlerFunc(handler.GetAll))
	router.Get("/ping", http.HandlerFunc(handler.GetPing))
	router.Get("/api/user/urls", http.HandlerFunc(handler.GetUserShortenedURLs))
	router.Get("/api/internal/stats", http.HandlerFunc(handler.GetStats))
	router.Delete("/api/user/urls", http.HandlerFunc(handler.DeleteUserShortenedURLs))

	return router
}
