package app

import (
	"sync"
	"testing"

	"package/main/internal/apexapi"
)

func TestModesStoreSetGet(t *testing.T) {
	var s modesStore
	if got := s.get().Pub.Current.Map; got != "" {
		t.Errorf("zero-value store should return empty map, got %q", got)
	}

	s.set(apexapi.Modes{Pub: apexapi.Maps{Current: apexapi.Map{Map: "Olympus"}}})
	if got := s.get().Pub.Current.Map; got != "Olympus" {
		t.Errorf("get after set = %q, want %q", got, "Olympus")
	}
}

// TestModesStoreConcurrent must be run with -race to be meaningful: it exercises
// concurrent writers and readers to prove the store is data-race free.
func TestModesStoreConcurrent(t *testing.T) {
	var s modesStore
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			s.set(apexapi.Modes{Pub: apexapi.Maps{Current: apexapi.Map{Map: "World's Edge"}}})
		}()
		go func() {
			defer wg.Done()
			_ = s.get().Pub.Current.Map
		}()
	}
	wg.Wait()
}

func TestChatSettingsStoreToggleEnabled(t *testing.T) {
	var s chatSettingsStore
	const chatID = int64(42)

	if got := s.enabled(chatID); got != false {
		t.Errorf("zero-value store should report disabled, got %v", got)
	}

	if got := s.toggle(chatID); got != true {
		t.Errorf("first toggle = %v, want true", got)
	}
	if got := s.enabled(chatID); got != true {
		t.Errorf("enabled after first toggle = %v, want true", got)
	}

	if got := s.toggle(chatID); got != false {
		t.Errorf("second toggle = %v, want false", got)
	}
	if got := s.enabled(chatID); got != false {
		t.Errorf("enabled after second toggle = %v, want false", got)
	}
}

// TestChatSettingsStoreConcurrent must be run with -race to be meaningful: it
// exercises concurrent togglers and readers to prove the store is data-race free.
func TestChatSettingsStoreConcurrent(t *testing.T) {
	var s chatSettingsStore
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			s.toggle(7)
		}()
		go func() {
			defer wg.Done()
			_ = s.enabled(7)
		}()
	}
	wg.Wait()
}
