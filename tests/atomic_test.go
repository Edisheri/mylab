// Тест для атомарного счетчика с эмуляцией задержек
package tests

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAtomicCounter(t *testing.T) {
	t.Parallel() // Параллельный запуск
	var count int32
	var wg sync.WaitGroup
	const goroutines = 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt32(&count, 1)
				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond) // Рандомные задержки
			}
		}()
	}
	wg.Wait()
	expected := int32(goroutines * 1000)
	if count != expected {
		t.Errorf("Expected %d, got %d", expected, count)
	}
}
