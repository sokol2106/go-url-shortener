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
	users sync.Map
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

	return token, err
}

func (ath *Authorization) IsUser(token string) (bool, error) {
	userID, err := ReadToken(token)
	if err != nil {
		return false, err
	}

	_, ok := ath.users.Load(userID)
	return ok, err
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
