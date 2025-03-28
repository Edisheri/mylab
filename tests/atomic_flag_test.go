package tests

import (
	"mylab/internal/atomic" 
	"sync"
	"testing"
)

func TestAtomicFlag(t *testing.T) {
	t.Parallel()
	flag := &atomic.AtomicFlag{}
	var wg sync.WaitGroup

	// 10 горутин устанавливают флаг
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			flag.Set()
		}()
	}

	// 5 горутин сбрасывают флаг
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			flag.Unset()
		}()
	}

	wg.Wait()
	if !flag.IsSet() {
		t.Error("Flag should be set, but it's not")
	}
}
