// Package cerrors предоставляет набор ошибок, используемых в приложении для обработки
// специфических условий и ситуаций, связанных с работой с короткими URL.
package cerrors

import "errors"

// ErrNewShortURL возникает, когда возникает конфликт при добавлении оригинального URL.
// Это может произойти, если уже существует запись с таким же оригинальным URL.
var ErrNewShortURL = errors.New("conflict add original URL")

// ErrGetShortURLDelete возникает, когда возникает конфликт при удалении оригинального URL.
// Это может произойти, если URL уже помечен как удалённый.
var ErrGetShortURLDelete = errors.New("conflict delete original URL")

// ErrGetShortURLNotFind возникает, когда оригинальный URL не может быть найден.
// Это может произойти, если URL не существует в базе данных.
var ErrGetShortURLNotFind = errors.New("conflict find original URL")
