package model

import "sync/atomic"

// SafeShortData представляет потокобезопасную оболочку для структуры ShortData.
// Доступ к данным осуществляется с использованием атомарных операций для обеспечения
// безопасного доступа из разных горутин.
type SafeShortData struct {
	value atomic.Value
}

// NewSafeShortData создает новый экземпляр SafeShortData и сохраняет начальные данные.
// Принимает объект ShortData в качестве параметра.
func NewSafeShortData(data ShortData) *SafeShortData {
	sd := &SafeShortData{}
	sd.value.Store(data)
	return sd
}

// Load возвращает текущее значение SafeShortData.
// Возвращаемое значение имеет тип ShortData.
func (sd *SafeShortData) Load() ShortData {
	return sd.value.Load().(ShortData)
}

// Store обновляет данные в SafeShortData новым значением ShortData.
// Операция является потокобезопасной.
func (sd *SafeShortData) Store(data ShortData) {
	sd.value.Store(data)
}
