package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"

type Token struct {
	jwt.RegisteredClaims
	UserID int
}

func NewToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Token{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},

		UserID: 1,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ReadToken(tknStrin string) (int, error) {
	token := &Token{}

	res, err := jwt.ParseWithClaims(tknStrin, token, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return -1, err
	}

	if !res.Valid {
		return -1, errors.New("Token is not valid")
	}

	return token.UserID, nil
}
