package progress

import (
	"sync"
	"testing"
)

func TestTrackerCounts(t *testing.T) {
	tracker := New()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		tracker.Print()
	}()

	for i := 0; i < 10; i++ {
		tracker.Increment()
		tracker.Done()
	}

	tracker.Wait()
	wg.Wait()

	if tracker.total != 10 {
		t.Fatalf("expected total 10, got %d", tracker.total)
	}

	if tracker.done != 10 {
		t.Fatalf("expected done 10, got %d", tracker.done)
	}
}
