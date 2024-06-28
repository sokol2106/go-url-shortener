package shorturl

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sokol2106/go-url-shortener/internal/gzip"
	"github.com/sokol2106/go-url-shortener/internal/logger"
	"io"
	"log"
	"net/http"
)

func NewShortURL(redirectURL string, fileStoragePath string) *ShortURL {
	s := new(ShortURL)
	s.redirectURL = redirectURL
	//s.shortDataList = su
	s.shortDataList.Init(fileStoragePath)
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

	res := s.shortDataList.AddURL(url)
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
	URL := s.shortDataList.GetURL(path)
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

func (s *ShortURL) PostJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Printf("ReadAll body error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		reqJS RequestJSON
		resJS ResponseJSON
	)

	err = json.Unmarshal(body, &reqJS)
	if err != nil {
		log.Printf("Unmarshal body error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := s.shortDataList.AddURL(reqJS.URL)
	resJS.Result = fmt.Sprintf("%s/%s", s.redirectURL, res)

	resBody, err := json.Marshal(resJS)
	if err != nil {
		log.Printf("Marshal body error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resBody)
}

func (s *ShortURL) Close() error {
	return s.shortDataList.Close()
}

func ShortRouter(sh *ShortURL) chi.Router {
	router := chi.NewRouter()

	// middleware
	router.Use(gzip.СompressionResponseRequest)
	router.Use(logger.LoggingResponseRequest)

	// router
	router.Post("/", http.HandlerFunc(sh.Post))
	router.Post("/api/shorten", http.HandlerFunc(sh.PostJSON))
	router.Get("/*", http.HandlerFunc(sh.GetAll))
	router.Get("/{id}", http.HandlerFunc(sh.Get))

	return router
}
