package main

import (
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"net/http"
)

func main() {

	servermux := http.NewServeMux()
	servermux.HandleFunc("/", shorturl.HanlerMain)

	err := http.ListenAndServe(":8080", servermux)
	if err != nil {

		panic(err)
	}

}
