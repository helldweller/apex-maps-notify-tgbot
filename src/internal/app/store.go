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

// chatSettingsStore tracks per-chat true-names toggle state in memory.
type chatSettingsStore struct {
	mu   sync.RWMutex
	data map[int64]bool
}

// toggle flips the true-names setting for chatID and returns the new state.
func (s *chatSettingsStore) toggle(chatID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data == nil {
		s.data = make(map[int64]bool)
	}
	s.data[chatID] = !s.data[chatID]
	return s.data[chatID]
}

// enabled reports whether true-names is on for chatID.
func (s *chatSettingsStore) enabled(chatID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[chatID]
}
