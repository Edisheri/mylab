package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// Счётчик, который будут увеличивать горутины
	var count int32

	// WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup
	wg.Add(2) // Заводим 2 горутины

	// Первая горутина
	go func() {
		defer wg.Done() // Уменьшаем счётчик wg при выходе
		for i := 0; i < 100000; i++ {
			// Атомарно увеличиваем count на 1
			atomic.AddInt32(&count, 1)
		}
	}()

	// Вторая горутина (аналогичная первой)
	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			atomic.AddInt32(&count, 1)
		}
	}()

	// Ждём, пока обе горутины завершат работу
	wg.Wait()

	// Выводим результат (всегда 200000)
	fmt.Println("Final count (atomic):", count)
}
