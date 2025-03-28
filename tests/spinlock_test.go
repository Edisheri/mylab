package tests

import (
	"mylab/internal/atomic" // Замените на реальный путь к вашему проекту
	"sync"
	"testing"
)

func TestSpinLock(t *testing.T) {
	t.Parallel() // Параллельный запуск [[4]]
	var spin atomic.SpinLock
	var count int
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 100000; i++ {
			spin.Lock()
			count++
			spin.Unlock()
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 100000; i++ {
			spin.Lock()
			count++
			spin.Unlock()
		}
		wg.Done()
	}()
	wg.Wait()

	if count != 200000 {
		t.Errorf("Expected 200000, got %d", count)
	}
}
