package main

import (
	"io"
	"net/http"
)

func hanlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		w.Write(body)
		w.Header().Set("Location", string(body))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		return
	}

	if r.Method == http.MethodGet {
		//w.Write([]byte("Temporary Redirect"))
		w.Header().Set("Location", "https://practicum.yandex.ru/")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(400)
	return

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", hanlerPost)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
