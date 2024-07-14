package shorturl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sokol2106/go-url-shortener/internal/cerrors"
	"github.com/sokol2106/go-url-shortener/internal/gzip"
	"github.com/sokol2106/go-url-shortener/internal/logger"
	"io"
	"log"
	"net/http"
)

func New(redirectURL string, strg StorageURL) *ShortURL {
	s := new(ShortURL)
	s.redirectURL = redirectURL
	s.storageURL = strg
	return s
}

func (s *ShortURL) createRedirectURL(url string) (string, error) {
	// НУЖНА ЛИ ПРОВЕРКА ВХОДНОГО URL !!!
	/*
		err := config.CheckURL(url)
		if err != nil {
			log.Printf("error CheckURL error: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return ""
		}
	*/
	res, err := s.storageURL.AddURL(url)
	return fmt.Sprintf("%s/%s", s.redirectURL, res), err
}

func (s *ShortURL) handlerError(err error) int {
	log.Printf("%s", err)
	if errors.Is(err, cerrors.ErrNewShortURL) {
		return http.StatusConflict
	}

	return http.StatusBadRequest
}

func (s *ShortURL) Post(w http.ResponseWriter, r *http.Request) {
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

	shortURL, err := s.createRedirectURL(string(body))
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

func (s *ShortURL) PostJSON(w http.ResponseWriter, r *http.Request) {
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

	resJS.Result, err = s.createRedirectURL(reqJS.URL)
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

func (s *ShortURL) PostBatch(w http.ResponseWriter, r *http.Request) {
	var (
		requestBatch  []RequestBatch
		responseBatch []ResponseBatch
		resBody       bytes.Buffer
		err           error
	)

	handlerStatus := http.StatusCreated

	if err := json.NewDecoder(r.Body).Decode(&requestBatch); err != nil {
		handlerStatus = s.handlerError(err)
		if handlerStatus == http.StatusBadRequest {
			w.WriteHeader(handlerStatus)
			return
		}
	}

	responseBatch, err = s.storageURL.AddBatch(requestBatch, s.redirectURL)
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

func (s *ShortURL) Get(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "id")
	URL := s.storageURL.GetURL(path)
	if URL != "" {
		w.Header().Set("Location", URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *ShortURL) GetAll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func (s *ShortURL) GetPing(w http.ResponseWriter, r *http.Request) {
	if s.storageURL == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := s.storageURL.PingContext()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ShortURL) Close() error {
	err := s.storageURL.Close()
	if err != nil {
		s.handlerError(err)
	}
	return err
}

func Router(sh *ShortURL) chi.Router {
	router := chi.NewRouter()

	// middleware
	router.Use(gzip.СompressionResponseRequest)
	router.Use(logger.LoggingResponseRequest)

	// router
	router.Post("/", http.HandlerFunc(sh.Post))
	router.Post("/api/shorten", http.HandlerFunc(sh.PostJSON))
	router.Post("/api/shorten/batch", http.HandlerFunc(sh.PostBatch))
	router.Get("/{id}", http.HandlerFunc(sh.Get))
	router.Get("/*", http.HandlerFunc(sh.GetAll))
	router.Get("/ping", http.HandlerFunc(sh.GetPing))

	return router
}
