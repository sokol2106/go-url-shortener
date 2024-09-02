package model

import "sync/atomic"

type SafeShortData struct {
	value atomic.Value
}

func NewSafeShortData(data ShortData) *SafeShortData {
	sd := &SafeShortData{}
	sd.value.Store(data)
	return sd
}

func (sd *SafeShortData) Load() ShortData {
	return sd.value.Load().(ShortData)
}

func (sd *SafeShortData) Store(data ShortData) {
	sd.value.Store(data)
}
