package shorturl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
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

func (s *ShortURL) createRedirectURL(url string) string {
	// НУЖНА ЛИ ПРОВЕРКА ВХОДНОГО URL !!!
	/*
		err := config.CheckURL(url)
		if err != nil {
			log.Printf("error CheckURL error: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return ""
		}
	*/
	res := s.storageURL.AddURL(url)
	return fmt.Sprintf("%s/%s", s.redirectURL, res)
}

func (s *ShortURL) handlerError(content string, err error) {
	log.Printf("Error %s: %s", content, err)
}

func (s *ShortURL) Post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		s.handlerError("read request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := s.createRedirectURL(string(body))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))

	if err != nil {
		s.handlerError("write response", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (s *ShortURL) Get(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "id")
	log.Printf("get get get : %s", path)
	URL := s.storageURL.GetURL(path)
	if URL != "" {
		w.Header().Set("Location", URL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *ShortURL) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("get ALL ")
	w.WriteHeader(http.StatusBadRequest)
}

func (s *ShortURL) GetPing(w http.ResponseWriter, r *http.Request) {
	if s.storageURL == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := s.storageURL.PingContext()
	if err != nil {
		s.handlerError("ping db", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ShortURL) PostJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		s.handlerError("read request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		reqJS RequestJSON
		resJS ResponseJSON
	)

	err = json.Unmarshal(body, &reqJS)
	if err != nil {
		s.handlerError("unmarshal body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resJS.Result = s.createRedirectURL(reqJS.URL)
	resBody, err := json.Marshal(resJS)
	if err != nil {
		s.handlerError("marshal body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resBody)
}

func (s *ShortURL) PostBatch(w http.ResponseWriter, r *http.Request) {
	var (
		requestBatch  []RequestBatch
		responseBatch []ResponseBatch
		resBody       bytes.Buffer
	)

	if err := json.NewDecoder(r.Body).Decode(&requestBatch); err != nil {
		s.handlerError("unmarshal body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseBatch = s.storageURL.AddBatch(requestBatch)
	if responseBatch != nil {
		err := json.NewEncoder(&resBody).Encode(responseBatch)
		if err != nil {
			s.handlerError("marshal body", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bodyB, err := io.ReadAll(&resBody)
		if err != nil {
			s.handlerError("read body", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		defer r.Body.Close()
		w.Write(bodyB)
	}

}

func (s *ShortURL) Close() error {
	err := s.storageURL.Close()
	if err != nil {
		s.handlerError("close storageURL", err)
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
