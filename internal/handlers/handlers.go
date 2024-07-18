package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sokol2106/go-url-shortener/internal/cerrors"
	"github.com/sokol2106/go-url-shortener/internal/gzip"
	"github.com/sokol2106/go-url-shortener/internal/handlers/token"
	"github.com/sokol2106/go-url-shortener/internal/logger"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"io"
	"log"
	"net/http"
)

type Handlers struct {
	srvShortURL *service.ShortURL
}

func NewHandlers(srv *service.ShortURL) *Handlers {
	return &Handlers{
		srvShortURL: srv,
	}
}

func (s *Handlers) handlerError(err error) int {
	statusCode := http.StatusBadRequest
	if errors.Is(err, cerrors.ErrNewShortURL) {
		statusCode = http.StatusConflict
	}

	log.Printf("error handling request: %v, status: %d", err, statusCode)
	return statusCode
}

func (s *Handlers) Post(w http.ResponseWriter, r *http.Request) {
	handlerStatus := http.StatusCreated
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	shortURL, err := s.srvShortURL.AddURL(string(body))
	if err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(handlerStatus)
	w.Write([]byte(shortURL))
}

func (s *Handlers) PostJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handlerStatus := http.StatusCreated
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	var (
		reqJS service.RequestJSON
		resJS service.ResponseJSON
	)

	err = json.Unmarshal(body, &reqJS)
	if err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	resJS.Result, err = s.srvShortURL.AddURL(reqJS.URL)
	if err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	resBody, err := json.Marshal(resJS)
	if err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(handlerStatus)
	w.Write(resBody)
}

func (s *Handlers) PostBatch(w http.ResponseWriter, r *http.Request) {
	var (
		requestBatch  []service.RequestBatch
		responseBatch []service.ResponseBatch
		resBody       bytes.Buffer
		err           error
	)

	handlerStatus := http.StatusCreated

	if err = json.NewDecoder(r.Body).Decode(&requestBatch); err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	responseBatch, err = s.srvShortURL.AddBatch(requestBatch)
	if err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	err = json.NewEncoder(&resBody).Encode(responseBatch)
	if err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	bodyB, err := io.ReadAll(&resBody)
	if err != nil {
		handlerStatus = s.handlerError(err)
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

func (s *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	path := chi.URLParam(r, "id")
	URL := s.srvShortURL.GetURL(ctx, path)
	if URL != "" {
		w.Header().Set("Location", URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Handlers) GetPing(w http.ResponseWriter, r *http.Request) {
	err := s.srvShortURL.PingContext()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Handlers) GetUserURL(w http.ResponseWriter, r *http.Request) {

}

func Router(sh *Handlers) chi.Router {
	router := chi.NewRouter()

	// middleware
	router.Use(gzip.СompressionResponseRequest)
	router.Use(logger.LoggingResponseRequest)
	router.Use(token.TokenResponseRequest)

	// router
	router.Post("/", http.HandlerFunc(sh.Post))
	router.Post("/api/shorten", http.HandlerFunc(sh.PostJSON))
	router.Post("/api/shorten/batch", http.HandlerFunc(sh.PostBatch))
	router.Get("/{id}", http.HandlerFunc(sh.Get))
	router.Get("/*", http.HandlerFunc(sh.GetAll))
	router.Get("/ping", http.HandlerFunc(sh.GetPing))
	router.Get("/api/user/urls", http.HandlerFunc(sh.GetUserURL))

	return router
}
