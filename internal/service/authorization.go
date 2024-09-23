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
	currentUserID string
	isNewUser     bool
}

func NewAuthorization() *Authorization {
	return &Authorization{}
}

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
	ath.isNewUser = false
	ath.currentUserID = userID
}

func (ath *Authorization) GetCurrentUserID() string {
	return ath.currentUserID
}

func (ath *Authorization) IsNewUser() bool {
	return ath.isNewUser
}

const tokenEXP = time.Hour * 24
const secretKey = "supersecret"

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
