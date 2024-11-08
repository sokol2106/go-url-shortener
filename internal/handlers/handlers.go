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
	srvShortURL   *service.ShortURL // Сервис сокращения URL
	trustedSubnet string            // Конфигурация сервера
}

// ResponseStats представляет структуру ответа GetStats
type ResponseStats struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}

// NewHandlers создает новый экземпляр Handlers с переданным сервисом сокращения URL.
func NewHandlers(srv *service.ShortURL, subnet string) *Handlers {
	return &Handlers{
		srvShortURL:   srv,
		trustedSubnet: subnet,
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

	shortURL, err := h.srvShortURL.AddOriginalURL(string(body), h.srvShortURL.GetAuthorization().GetCurrentUserID())

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

	resBody, err := h.srvShortURL.AddOriginalURLJSON(body, h.srvShortURL.GetAuthorization().GetCurrentUserID())

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

	responseBatch, err = h.srvShortURL.AddOriginalURLBatch(requestBatch, h.srvShortURL.GetAuthorization().GetCurrentUserID())

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
	if h.srvShortURL.GetAuthorization().IsNewUser() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	res, err := h.srvShortURL.GetUserShortenedURLs(ctx, h.srvShortURL.GetAuthorization().GetCurrentUserID())
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
	if h.srvShortURL.GetAuthorization().IsNewUser() {
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

	h.srvShortURL.DeleteOriginalURLs(ctx, h.srvShortURL.GetAuthorization().GetCurrentUserID(), request)

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

	res.Users = h.srvShortURL.GetAuthorization().GetUsers()
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
func (h *Handlers) Router() chi.Router {
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.СompressionResponseRequest)
	router.Use(middleware.LoggingResponseRequest)
	router.Use(func(handlerF http.Handler) http.Handler {
		return middleware.TokenResponseRequest(h.srvShortURL.GetAuthorization(), handlerF)
	})

	// router
	router.Post("/", http.HandlerFunc(h.Post))
	router.Post("/api/shorten", http.HandlerFunc(h.PostJSON))
	router.Post("/api/shorten/batch", http.HandlerFunc(h.PostBatch))
	router.Get("/{id}", http.HandlerFunc(h.Get))
	router.Get("/*", http.HandlerFunc(h.GetAll))
	router.Get("/ping", http.HandlerFunc(h.GetPing))
	router.Get("/api/user/urls", http.HandlerFunc(h.GetUserShortenedURLs))
	router.Get("/api/internal/stats", http.HandlerFunc(h.GetStats))
	router.Delete("/api/user/urls", http.HandlerFunc(h.DeleteUserShortenedURLs))

	return router
}
