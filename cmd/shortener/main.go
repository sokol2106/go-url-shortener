package main

import "net/http"

func hanlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Write([]byte("http://localhost:8080/EwHXdJfB"))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(201)
		return
	}

	if r.Method == http.MethodGet {
		w.Write([]byte("Temporary Redirect"))
		w.Header().Set("Location", "https://practicum.yandex.ru/")
		w.WriteHeader(307)
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
