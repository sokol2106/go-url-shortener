package shorturl

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	storage "github.com/sokol2106/go-url-shortener/internal/storage"
	"io"
	"math/rand"
	"net/http"
	"net/url"
)

var tableshortdata = make(map[string]*storage.Shortdata)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func HanlerMain(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = checkURL(string(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hash := md5.Sum(body)
		thash := hex.EncodeToString(hash[:])

		tshdata, exist := tableshortdata[thash]
		if !exist {
			tshdata = storage.NewShortdata(string(body), "/"+randText(8))
			tableshortdata[thash] = tshdata
		}

		//w.Header().Set("Location", thash)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(w, "http://localhost:8080%s", tshdata.Short())

		return
	}

	if r.Method == http.MethodGet {
		path := r.URL.Path

		for _, value := range tableshortdata {
			if path == value.Short() {
				w.Header().Set("Location", value.URL())
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}
	}

	w.WriteHeader(http.StatusBadRequest)
}

func checkURL(body string) error {
	urlParse, err := url.Parse(body)
	if err != nil {
		return err
	}

	if urlParse.Scheme != "http" && urlParse.Scheme != "https" || urlParse.Host == "" {

		return errors.New("invalid url")
	}

	return nil
}

func randText(size int) string {
	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
