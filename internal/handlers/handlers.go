package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sokol2106/go-url-shortener/internal/cerrors"
	"github.com/sokol2106/go-url-shortener/internal/gzip"
	"github.com/sokol2106/go-url-shortener/internal/logger"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"io"
	"log"
	"net/http"
)

type Handlers struct {
	srvShortURL      *service.ShortURL
	srvAuthorization *service.Authorization
}

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

func NewHandlers(srv *service.ShortURL) *Handlers {
	return &Handlers{
		srvShortURL:      srv,
		srvAuthorization: service.NewAuthorization(),
	}
}

func (h *Handlers) handlerError(err error) int {
	statusCode := http.StatusBadRequest
	if errors.Is(err, cerrors.ErrNewShortURL) {
		statusCode = http.StatusConflict
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
	URL := h.srvShortURL.GetOriginalURL(ctx, path)
	if URL != "" {
		w.Header().Set("Location", URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
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

	log.Printf("err AAA request: len: %d", len(res))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (h *Handlers) TokenResponseRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user")
		// не существует или она не проходит проверку подлинности
		if err != nil {
			h.handlerError(err)
			tkn, err := h.srvAuthorization.NewUserToken()
			if err != nil {
				h.handlerError(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			newCookie := http.Cookie{Name: "user", Value: tkn}
			http.SetCookie(w, &newCookie)
		} else {
			// Без ошибок
			userID, err := h.srvAuthorization.GetUserID(cookie.Value)
			if err != nil {
				h.handlerError(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			isUser := h.srvAuthorization.IsUser(userID)
			if !isUser {
				h.handlerError(err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.srvAuthorization.SetCurrentUserID(userID)
			http.SetCookie(w, cookie)
		}

		handler.ServeHTTP(w, r)
	})

}

func Router(handler *Handlers) chi.Router {
	router := chi.NewRouter()

	// middleware
	router.Use(gzip.СompressionResponseRequest)
	router.Use(logger.LoggingResponseRequest)
	router.Use(handler.TokenResponseRequest)

	// router
	router.Post("/", http.HandlerFunc(handler.Post))
	router.Post("/api/shorten", http.HandlerFunc(handler.PostJSON))
	router.Post("/api/shorten/batch", http.HandlerFunc(handler.PostBatch))
	router.Get("/{id}", http.HandlerFunc(handler.Get))
	router.Get("/*", http.HandlerFunc(handler.GetAll))
	router.Get("/ping", http.HandlerFunc(handler.GetPing))
	router.Get("/api/user/urls", http.HandlerFunc(handler.GetUserShortenedURLs))

	return router
}
