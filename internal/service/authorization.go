package service

import "sync"

type Authorization struct {
	users sync.Map
}

func NewAuthorization() *Authorization {
	return &Authorization{}
}

func (ath *Authorization) NewUser(login string) {

}

func IsUser(login string) bool {
	return false
}
