package progress

import (
	"fmt"
	"sync"
	"time"
)

type Tracker struct {
	mu    sync.Mutex
	total int
	done  int
	wg    sync.WaitGroup
}

func New() *Tracker {
	t := &Tracker{}
	t.wg.Add(1)
	return t
}

func (t *Tracker) Increment() {
	t.mu.Lock()
	t.total++
	t.mu.Unlock()
}

func (t *Tracker) Done() {
	t.mu.Lock()
	t.done++
	t.mu.Unlock()
}

func (t *Tracker) Print() {
	defer t.wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		t.mu.Lock()
		fmt.Printf("\rProcessed %d / %d", t.done, t.total)
		if t.done == t.total {
			t.mu.Unlock()
			return
		}
		t.mu.Unlock()
	}
}

func (t *Tracker) Wait() {
	t.wg.Wait()
}
