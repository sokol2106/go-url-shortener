package shorturl

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sokol2106/go-url-shortener/internal/server"
	storage "github.com/sokol2106/go-url-shortener/internal/storage"
	"io"
	"math/rand"
	"net/http"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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

func randText(size int) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (s *ShortURL) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = server.CheckURL(string(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hash := md5.Sum(body)
		thash := hex.EncodeToString(hash[:])

		tshdata, exist := s.tableshortdata[thash]
		if !exist {
			tshdata = storage.NewShortdata(string(body), randText(8))
			s.tableshortdata[thash] = tshdata
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, "%s/%s", s.url, tshdata.Short())

		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (s *ShortURL) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		path := chi.URLParam(r, "id")

		for _, value := range s.tableshortdata {
			if path == value.Short() {
				w.Header().Set("Location", value.URL())
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}
	}

	w.WriteHeader(http.StatusBadRequest)
}

func ShortRouter(url string) chi.Router {
	router := chi.NewRouter()
	sh := NewShortURL(url)
	router.Post("/", sh.Post)
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadRequest) })
	router.Get("/{id}", sh.Get)

	return router
}
