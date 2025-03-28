# Atomic Variables in Go: Исследование и примеры

## Содержание
- [Введение](#введение)
- [Критическая ситуация: Гонки данных](#критическая-ситуация-гонки-данных)
- [Решение с атомарными операциями](#решение-с-атомарными-операциями)
- [Паттерн: Спин-лок](#паттерн-спин-лок)
- [Автоматические тесты](#автоматические-тесты)
- [Дополнительные задачи](#дополнительные-задачи)
- [Ссылки](#ссылки)
- [Видео](#видео)

---

## Введение
Атомарные переменные — это ключевой инструмент для безопасной работы с общими данными в многопоточных приложениях на Go. Они обеспечивают атомарность операций (неделимость), предотвращая гонки данных (race conditions). В этом проекте:
- Рассматривается проблема гонок данных на примере счетчика.
- Приводится решение с использованием пакета `sync/atomic`.
- Реализованы тесты с параллельным выполнением и эмуляцией задержек.
- Предложены дополнительные задачи на concurrency.

---

## Критическая ситуация: Гонки данных
### Описание проблемы
При одновременном доступе нескольких горутин к общей переменной возникает состояние гонки. Например, операция `count++` неатомарна и состоит из трех шагов:
1. Чтение значения.
2. Увеличение на 1.
3. Запись результата.

### Пример кода с ошибкой
```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var count int
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        for i := 0; i < 100000; i++ {
            count++
        }
        wg.Done()
    }()
    go func() {
        for i := 0; i < 100000; i++ {
            count++
        }
        wg.Done()
    }()
    wg.Wait()
    fmt.Println("Final count:", count) // Результат < 200000
}
```


### Запуск:
````
go run examples/race_condition/main.go
````

### Решение с атомарными операциями
Используем sync/atomic для атомарного инкремента:

````
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var count int32
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        for i := 0; i < 100000; i++ {
            atomic.AddInt32(&count, 1)
        }
        wg.Done()
    }()
    go func() {
        for i := 0; i < 100000; i++ {
            atomic.AddInt32(&count, 1)
        }
        wg.Done()
    }()
    wg.Wait()
    fmt.Println("Final count:", count) // Всегда 200000
}
````

### Запуск:
````
go run examples/atomic_counter/main.go
````

### Паттерн: Спин-лок
Реализация спин-лока через атомарные операции:

````
type SpinLock struct {
    locked int32
}

func (sl *SpinLock) Lock() {
    for !atomic.CompareAndSwapInt32(&sl.locked, 0, 1) {
        // Спин-ожидание
    }
}

func (sl *SpinLock) Unlock() {
    atomic.StoreInt32(&sl.locked, 0)
}
````
### Пример использования:
````
var spin SpinLock
var shared int

func increment() {
    spin.Lock()
    shared++
    spin.Unlock()
}
````

## Автоматические тесты
### Тест для атомарного счетчика

````
func TestAtomicCounter(t *testing.T) {
    var count int32
    var wg sync.WaitGroup
    const goroutines = 100

    wg.Add(goroutines)
    for i := 0; i < goroutines; i++ {
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                atomic.AddInt32(&count, 1)
                time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
            }
        }()
    }
    wg.Wait()
    expected := int32(goroutines * 1000)
    if count != expected {
        t.Errorf("Expected %d, got %d", expected, count)
    }
}
````
### Запуск тестов
````
go test -race -v ./tests/...
````

### Дополнительные задачи
### Задача 1: Безопасный счетчик
### Условие:
Реализуйте счетчик, который можно безопасно увеличивать/уменьшать из нескольких горутин.

### Задача 2: Безопасный счетчик
### Условие: Реализуйте флаг, который можно атомарно устанавливать и проверять.

## Ссылки на академические ресурсы
1. **Книга:** ["Learn Concurrent Programming with Go" by James Cutajar](https://www.amazon.com/Learn-Concurrent-Programming-Golang-practical/dp/1803245185)
    - Подробное введение в конкурентность в Go, включая атомарные переменные.

2. **Статья:** ["Atomic Operations in Golang - atomic package"](https://go.dev/doc/atomic)
    - Официальная документация Go, объясняющая основы атомарных операций.

3. **Сайт:** [GeeksforGeeks - "Atomic Variable in Golang"](https://www.geeksforgeeks.org/atomic-variables-in-golang/)
    - Четкое объяснение атомарных переменных с примерами.

4. **Блог на Medium:** ["The Go 1.19 Atomic Wrappers and why to use them" by Ralph Caraveo](https://medium.com/@rsc/the-go-1-19-atomic-wrappers-and-why-to-use-them-9b0d1a3d8e3c)
    - Обсуждение новых атомарных оберток в Go 1.19.

5. **Stack Overflow:** ["Is variable assignment atomic in go?"](https://stackoverflow.com/questions/36609303/is-variable-assignment-atomic-in-go)
    - Обсуждение атомарности присваивания переменных в Go.

---

## Ссылки на статьи в интернете
1. [Go by Example: Atomic Counters](https://gobyexample.com/atomic-counters)
2. [Atomic Operations Provided in The sync/atomic Standard Package](https://go101.org/article/sync-atomic.html)
3. [An Introduction to Atomic Variables in Java](https://www.baeldung.com/java-atomic-variables) (концепция релевантна для Go)
4. [Atomic Variable in Golang](https://www.tutorialspoint.com/atomic-variables-in-golang)
5. [Atomic. C++ and Go.. First article will be about comparing…](https://medium.com/@blanchon.vincent/atomic-c-and-go-3e9a2b4d0e4e)
6. [Is variable assignment atomic in go?](https://stackoverflow.com/questions/36609303/is-variable-assignment-atomic-in-go)
7. [Go Atomic Variable](https://www.tpointtech.com/golang/atomic-variables-in-golang/)
8. [Atomic Operations in Golang](https://www.geeksforgeeks.org/atomic-operations-in-golang/)
9. [The Go 1.19 Atomic Wrappers and why to use them](https://medium.com/@rsc/the-go-1-19-atomic-wrappers-and-why-to-use-them-9b0d1a3d8e3c)

---

## Видео с объяснением
- [Google I/O 2012 - Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
    - Обзор паттернов конкурентности в Go, включая атомарные операции.

---

## Ключевые цитирования
1. **James Cutajar** в книге ["Learn Concurrent Programming with Go"](https://www.amazon.com/Learn-Concurrent-Programming-Golang-practical/dp/1803245185).
2. [Atomic Operations in Golang - atomic package](https://go.dev/doc/atomic) (Golang Docs).
3. [Atomic Variable in Golang](https://www.geeksforgeeks.org/atomic-variables-in-golang/) (GeeksforGeeks).
4. [The Go 1.19 Atomic Wrappers](https://medium.com/@rsc/the-go-1-19-atomic-wrappers-and-why-to-use-them-9b0d1a3d8e3c) (Medium).
5. [Is variable assignment atomic in go?](https://stackoverflow.com/questions/36609303/) (Stack Overflow).
6. [sync/atomic](https://pkg.go.dev/sync/atomic) (Go Packages).
7. [Go by Example: Atomic Counters](https://gobyexample.com/atomic-counters).
8. [Atomic Operations in The sync/atomic Package](https://go101.org/article/sync-atomic.html) (Go 101).
9. [Atomic Variables in Java](https://www.baeldung.com/java-atomic-variables) (Baeldung).
10. [Go Atomic Variable](https://www.tpointtech.com/golang/atomic-variables-in-golang/) (Tpoint Tech).  