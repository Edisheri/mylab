package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int // Общая переменная
	var wg sync.WaitGroup
	wg.Add(2) // Будет 2 горутины

	// Горутина 1
	go func() {
		for i := 0; i < 100000; i++ {
			count++ //
		}
		wg.Done()
	}()

	// Горутина 2 (аналогичная)
	go func() {
		for i := 0; i < 100000; i++ {
			count++
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Final count:", count) // Результат < 200000
}
