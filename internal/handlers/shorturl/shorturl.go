package shorturl

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/logger"
	storage "github.com/sokol2106/go-url-shortener/internal/storage"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func randText(size int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

type ShortURL struct {
	url            string
	tableshortdata map[string]*storage.Shortdata
}

func NewShortURL(u string) *ShortURL {
	return &ShortURL{
		url:            u,
		tableshortdata: make(map[string]*storage.Shortdata),
	}
}

func (s *ShortURL) Post() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = config.CheckURL(string(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*
			hash := sha256.Sum256(body)
			thash := hex.EncodeToString(hash[:])
			tshdata, exist := s.tableshortdata[thash]
			if !exist {
				tshdata = storage.NewShortdata(string(body), randText(8))
				s.tableshortdata[thash] = tshdata
			}
		*/

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, s.addURL(string(body)))
		//_, _ = fmt.Fprintf(w, "%s/%s", s.url, tshdata.Short())
	}

	return http.HandlerFunc(fn)
}

func (s *ShortURL) Get() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/*
			path := chi.URLParam(r, "id")
			for _, value := range s.tableshortdata {
				if path == value.Short() {
					w.Header().Set("Location", value.URL())
					w.WriteHeader(http.StatusTemporaryRedirect)
					return
				}
			}

		*/

		path := chi.URLParam(r, "id")
		URL := s.getURL(path)
		if URL != "" {
			w.Header().Set("Location", URL)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}

	return http.HandlerFunc(fn)
}

func (s *ShortURL) GetAll() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}
	return http.HandlerFunc(fn)
}

func (s *ShortURL) PostJSON() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	return http.HandlerFunc(fn)
}

func (s *ShortURL) addURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	thash := hex.EncodeToString(hash[:])
	tshdata, exist := s.tableshortdata[thash]
	if !exist {
		tshdata = storage.NewShortdata(url, randText(8))
		s.tableshortdata[thash] = tshdata
		return fmt.Sprintf("%s/%s", s.url, tshdata.Short())
	}

	return ""
}

func (s *ShortURL) getURL(shURL string) string {
	for _, value := range s.tableshortdata {
		if shURL == value.Short() {
			return value.URL()
		}
	}
	return ""
}

func ShortRouter(url string) chi.Router {
	router := chi.NewRouter()
	sh := NewShortURL(url)
	router.Post("/", logger.LoggingResponseRequest(sh.Post()))
	router.Post("/api/shorten", logger.LoggingResponseRequest(sh.PostJSON()))
	router.Get("/*", logger.LoggingResponseRequest(sh.GetAll()))
	router.Get("/{id}", logger.LoggingResponseRequest(sh.Get()))

	return router
}
