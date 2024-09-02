package middleware

import (
	"github.com/sokol2106/go-url-shortener/internal/service"
	"log"
	"net/http"
)

type Token struct {
	srvAuthorization *service.Authorization
}

func NewToken() *Token {
	return &Token{service.NewAuthorization()}
}

func (t *Token) GetAuthorization() *service.Authorization {
	return t.srvAuthorization
}

func (t *Token) TokenResponseRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user")
		// не существует или она не проходит проверку подлинности
		if err != nil {
			tkn, err := t.srvAuthorization.NewUserToken()
			if err != nil {
				log.Printf("error handling request: %v, status: %d", err, http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			newCookie := http.Cookie{Name: "user", Value: tkn}
			http.SetCookie(w, &newCookie)
		} else {
			// Без ошибок
			userID, err := t.srvAuthorization.GetUserID(cookie.Value)
			if err != nil {
				log.Printf("error handling request: %v, status: %d", err, http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			isUser := t.srvAuthorization.IsUser(userID)
			if !isUser {
				log.Printf("error handling request: %v, status: %d", err, http.StatusInternalServerError)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			t.srvAuthorization.SetCurrentUserID(userID)
			http.SetCookie(w, cookie)
		}

		handler.ServeHTTP(w, r)
	})
}
