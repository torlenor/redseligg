package utils

import (
	"sync"
)

// IDProvider takes care of returning a fresh ID every time it is asked for
// It is thread safe
type IDProvider struct {
	lastID int
	mutex  sync.Mutex
}

// Get an unique ID
func (i *IDProvider) Get() int {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.lastID++

	return i.lastID
}
