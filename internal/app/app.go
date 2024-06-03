package app

import (
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/server"
)

func Run() {
	/*servermux := http.NewServeMux()
	servermux.HandleFunc("/", shorturl.HanlerMain)

	err := http.ListenAndServe(":8080", servermux)
	if err != nil {

		panic(err)
	}

	*/
	
	srv := server.NewServer(shorturl.ShortRouter())
	err := srv.Start()
	if err != nil {

		panic(err)
	}
}
