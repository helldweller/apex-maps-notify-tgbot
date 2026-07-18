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
