package main

import (
	"math/rand/v2"
	"sync"
)

type Scribble struct {
	Date [3]uint `json:"date"` // yyyy mm dd
	Text string  `json:"text"`
}

var scribbleCache = make([]Scribble, 0)
var cacheMutex sync.RWMutex

func randScribble() (*Scribble, bool) {
	if len(scribbleCache) == 0 {
		return nil, false
	}

	cacheMutex.RLock()

	scribble := scribbleCache[rand.IntN(len(scribbleCache))]

	cacheMutex.RUnlock()

	return &scribble, true
}

func addScribble(scribble Scribble) {
	cacheMutex.Lock()

	scribbleCache = append(scribbleCache, scribble)

	cacheMutex.Unlock()
}
