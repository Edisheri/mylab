package tests

import (
	"mylab/internal/atomic" // Замените на реальный путь
	"sync"
	"testing"
)

func TestSafeCounter(t *testing.T) {
	t.Parallel()
	counter := &atomic.SafeCounter{}
	var wg sync.WaitGroup

	// 100 горутин увеличивают счетчик
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	// 50 горутин уменьшают счетчик
	wg.Add(50)
	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()
			counter.Decrement()
		}()
	}

	wg.Wait()
	expected := int32(100 - 50)
	if counter.Value() != expected {
		t.Errorf("Expected %d, got %d", expected, counter.Value())
	}
}
