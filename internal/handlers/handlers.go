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
	srvShortURL *service.ShortURL
	srvAuth     *service.Authorization
}

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

func NewHandlers(srv *service.ShortURL) *Handlers {
	return &Handlers{
		srvShortURL: srv,
		srvAuth:     service.NewAuthorization(),
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
		if s.handlerError(err) == http.StatusBadRequest {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	user := "123456"
	shortURL, err := s.srvShortURL.AddURL(string(body), user)

	if err != nil {
		if s.handlerError(err) == http.StatusBadRequest {
			w.WriteHeader(http.StatusBadRequest)
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
		reqJS RequestJSON
		resJS ResponseJSON
	)

	err = json.Unmarshal(body, &reqJS)
	if err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	user := "123456"
	resJS.Result, err = s.srvShortURL.AddURL(reqJS.URL, user)

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

	userID := "123456"
	responseBatch, err = s.srvShortURL.AddBatch(requestBatch, userID)
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

func (s *Handlers) TokenResponseRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*cookie, err := r.Cookie("user")
		// Ошибка по куке
		if err != nil {
			s.handlerError(err)
			tkn, err := s.srvAuth.NewUserToken()
			if err != nil {
				s.handlerError(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			newCookie := http.Cookie{Name: "user", Value: tkn}
			http.SetCookie(w, &newCookie)
		} else {
			// Без ошибок
			isUser, err := s.srvAuth.IsUser(cookie.Value)
			if err != nil {
				s.handlerError(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !isUser {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.SetCookie(w, cookie)
		}*/

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
	router.Get("/api/user/urls", http.HandlerFunc(handler.GetUserURL))

	return router
}
