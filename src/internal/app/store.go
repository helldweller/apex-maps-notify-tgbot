package app

import (
	"sync"

	"package/main/internal/apexapi"
)

// modesStore holds the latest map rotation and guards it against concurrent
// access by the updater goroutine (writer) and the Telegram handler (reader).
type modesStore struct {
	mu   sync.RWMutex
	data apexapi.Modes
}

// set atomically replaces the stored rotation.
func (s *modesStore) set(m apexapi.Modes) {
	s.mu.Lock()
	s.data = m
	s.mu.Unlock()
}

// get returns a snapshot of the stored rotation. apexapi.Modes is a plain value
// struct, so the returned copy is safe to read without further locking.
func (s *modesStore) get() apexapi.Modes {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data
}
