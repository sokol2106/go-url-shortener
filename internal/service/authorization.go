// Package service предоставляет сервисы сисетмы
// Сервис управления авторизацией пользователей с использованием JWT (JSON Web Token) и хранения информации о пользователях в системе.

package service

import (
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"math/big"
	"sync"
	"sync/atomic"
	"time"
)

// Authorization хранит информацию о пользователях системы и их текущем состоянии авторизации.
type Authorization struct {
	users         sync.Map // пользователи системы, где ключ — UserID, а значение — информация о пользователе
	currentUserID string   // ID текущего авторизованного пользователя
	isNewUser     bool     // флаг, указывающий, что это новый пользователь
	countUsers    int64    // количество пользователей
}

// NewAuthorization создает новый объект Authorization и возвращает указатель на него.
//
// Пример использования:
//
//	authorization := service.NewAuthorization()
func NewAuthorization() *Authorization {
	return &Authorization{
		countUsers: 0,
	}
}

// NewUserToken создает новый JWT-токен для нового пользователя.
// Токен содержит уникальный UserID, который также сохраняется в объекте Authorization.
// Возвращает строку с токеном или ошибку.
//
// Пример использования:
//
//	token, err := authorization.NewUserToken()
//	if err != nil {
//		log.Fatalf("Ошибка при создании токена: %v", err)
//	}
func (ath *Authorization) NewUserToken() (string, error) {
	userID, err := rand.Int(rand.Reader, big.NewInt(30000))
	user := "user1"
	if err != nil {
		return "", err
	}

	ath.users.Store(userID.String(), user)
	token, err := NewToken(userID.String())

	ath.currentUserID = userID.String()
	ath.isNewUser = true

	atomic.AddInt64(&ath.countUsers, 1)
	return token, err
}

// IsUser проверяет, зарегистрирован ли пользователь с данным userID.
// Возвращает true, если пользователь найден.
func (ath *Authorization) IsUser(userID string) bool {
	_, ok := ath.users.Load(userID)
	return ok
}

// GetUserID извлекает UserID из JWT-токена.
// Возвращает UserID и ошибку, если токен не валиден.
func (ath *Authorization) GetUserID(token string) (string, error) {
	userID, err := ReadToken(token)
	if err != nil {
		return "", err
	}

	return userID, err
}

// SetCurrentUserID устанавливает UserID как текущего авторизованного пользователя.
func (ath *Authorization) SetCurrentUserID(userID string) {
	ath.isNewUser = false
	ath.currentUserID = userID
}

// GetUsers возвращает количество пользователей
func (ath *Authorization) GetUsers() int {
	return int(atomic.LoadInt64(&ath.countUsers))
}

// GetCurrentUserID возвращает UserID текущего авторизованного пользователя.
func (ath *Authorization) GetCurrentUserID() string {
	return ath.currentUserID
}

// IsNewUser возвращает true, если текущий пользователь является новым.
func (ath *Authorization) IsNewUser() bool {
	return ath.isNewUser
}

// Длительность жизни JWT-токена (24 часа).
const tokenEXP = time.Hour * 24

// Секретный ключ для подписи JWT-токенов.
const secretKey = "supersecret"

// Token описывает структуру JWT-токена с полем UserID.
type Token struct {
	jwt.RegisteredClaims
	UserID string
}

// NewToken создает и возвращает новый JWT-токен с указанным userID.
// Возвращает строку с токеном или ошибку.
func NewToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Token{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenEXP)),
		},

		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ReadToken проверяет валидность JWT-токена и возвращает UserID из токена.
// Возвращает UserID и ошибку, если токен недействителен.
func ReadToken(cookValue string) (string, error) {
	token := &Token{}

	res, err := jwt.ParseWithClaims(cookValue, token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !res.Valid {
		return "", errors.New("Token is not valid")
	}

	return token.UserID, nil
}
