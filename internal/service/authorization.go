package service

import (
	"crypto/rand"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"math/big"
	"sync"
	"time"
)

type Authorization struct {
	users         sync.Map
	currentUserId string
}

func NewAuthorization() *Authorization {
	return &Authorization{}
}

func (ath *Authorization) NewUserToken() (string, error) {
	userID, err := rand.Int(rand.Reader, big.NewInt(15))
	user := "user1"
	if err != nil {
		return "", err
	}

	ath.users.Store(userID.String(), user)
	token, err := NewToken(userID.String())

	ath.currentUserId = userID.String()
	return token, err
}

func (ath *Authorization) IsUser(userID string) bool {
	_, ok := ath.users.Load(userID)
	return ok
}

func (ath *Authorization) GetUserID(token string) (string, error) {
	userID, err := ReadToken(token)
	if err != nil {
		return "", err
	}

	return userID, err
}

func (ath *Authorization) SetCurrentUserID(userID string) {
	ath.currentUserId = userID
}

func (ath *Authorization) GetCurrentUserID() string {
	return ath.currentUserId
}

const tokenEXP = time.Hour * 3
const secretKey = "supersecretkey"

type Token struct {
	jwt.RegisteredClaims
	UserID string
}

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
