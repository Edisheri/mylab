package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Спин-блокировка (SpinLock)
type SpinLock struct {
	locked int32 // 0 — свободен, 1 — занят
}

// Lock захватывает блокировку
func (s *SpinLock) Lock() {
	// Пытаемся атомарно изменить locked с 0 на 1
	for !atomic.CompareAndSwapInt32(&s.locked, 0, 1) {
		// Если не получилось — повторяем попытку
	}
}

// Unlock освобождает блокировку
func (s *SpinLock) Unlock() {
	atomic.StoreInt32(&s.locked, 0) // Атомарно сбрасываем в 0
}

func main() {
	var spin SpinLock
	var count int
	var wg sync.WaitGroup
	wg.Add(2)

	// Горутина 1
	go func() {
		for i := 0; i < 100000; i++ {
			spin.Lock()   // Захватываем блокировку
			count++       // Изменяем общую переменную
			spin.Unlock() // Освобождаем блокировку
		}
		wg.Done()
	}()

	// Горутина 2 (аналогичная)
	go func() {
		for i := 0; i < 100000; i++ {
			spin.Lock()
			count++
			spin.Unlock()
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Final count (spinlock):", count) // Всегда 200000
}
