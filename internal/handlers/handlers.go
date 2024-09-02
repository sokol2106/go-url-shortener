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
	"net/http"
)

type Handlers struct {
	srvShortURL      *service.ShortURL
	token            *middleware.Token
	srvAuthorization *service.Authorization
}

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

func NewHandlers(srv *service.ShortURL) *Handlers {
	t := middleware.NewToken()
	a := t.GetAuthorization()
	return &Handlers{
		srvShortURL:      srv,
		token:            t,
		srvAuthorization: a,
	}
}

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

func (h *Handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func (h *Handlers) GetPing(w http.ResponseWriter, r *http.Request) {
	err := h.srvShortURL.PingContext()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

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
	router.Delete("/api/user/urls", http.HandlerFunc(handler.DeleteUserShortenedURLs))

	return router
}
