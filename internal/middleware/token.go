package middleware

import (
	"github.com/sokol2106/go-url-shortener/internal/service"
	"log"
	"net/http"
)

// TokenResponseRequest является middleware-обработчиком, который проверяет наличие куки с токеном "user".
// Если куки не существует или токен недействителен, создает новый токен и устанавливает его в куки.
// Если токен существует и действителен, проверяет пользователя и продолжает выполнение запроса.
// В случае ошибки возвращает соответствующий HTTP-статус.
func TokenResponseRequest(srvAuthorization *service.Authorization, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user")
		// не существует или она не проходит проверку подлинности
		if err != nil {
			tkn, err := srvAuthorization.NewUserToken()
			if err != nil {
				log.Printf("error handling request: %v, status: %d", err, http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			newCookie := http.Cookie{Name: "user", Value: tkn}
			http.SetCookie(w, &newCookie)
		} else {
			// Без ошибок
			userID, err := srvAuthorization.GetUserID(cookie.Value)
			if err != nil {
				log.Printf("error handling request: %v, status: %d", err, http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			isUser := srvAuthorization.IsUser(userID)
			if !isUser {
				log.Printf("error handling request: %v, status: %d", err, http.StatusInternalServerError)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			srvAuthorization.SetCurrentUserID(userID)
			http.SetCookie(w, cookie)
		}

		handler.ServeHTTP(w, r)
	})
}
