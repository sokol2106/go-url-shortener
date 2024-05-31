package main

import (
	shorturl "github.com/sokol2106/go-url-shortener/internal/app"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", shorturl.HanlerMain)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {

		panic(err)
	}

}
