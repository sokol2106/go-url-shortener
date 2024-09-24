package middleware

import (
	"github.com/sokol2106/go-url-shortener/internal/service"
	"log"
	"net/http"
)

// Token представляет структуру для работы с сервисом авторизации Authorization.
// srvAuthorization - объект сервиса авторизации, который используется для создания и проверки токенов.
type Token struct {
	srvAuthorization *service.Authorization
}

// NewToken инициализирует объект Token и service.Authorization.
// Возвращает указатель на Token.
func NewToken() *Token {
	return &Token{service.NewAuthorization()}
}

// GetAuthorization возвращает указатель на объект сервиса авторизации service.Authorization.
func (t *Token) GetAuthorization() *service.Authorization {
	return t.srvAuthorization
}

// TokenResponseRequest является middleware-обработчиком, который проверяет наличие куки с токеном "user".
// Если куки не существует или токен недействителен, создает новый токен и устанавливает его в куки.
// Если токен существует и действителен, проверяет пользователя и продолжает выполнение запроса.
// В случае ошибки возвращает соответствующий HTTP-статус.
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
