# Atomic Variables in Go: Исследование и примеры

## Содержание
- [Введение](#введение)
- [Что такое атомарные переменные?](#что-такое-атомарные-переменные)
- [CAS-операции (Compare-And-Swap)](#cas-операции)
- [Критическая ситуация: Гонки данных](#критическая-ситуация-гонки-данных)
- [Решение с атомарными операциями](#решение-с-атомарными-операциями)
- [Паттерн: Спин-лок](#паттерн-спин-лок)
- [Автоматические тесты](#автоматические-тесты)
- [Дополнительные задачи](#дополнительные-задачи)
- [Ссылки](#ссылки)
- [Видео](#видео)

---

## Введение
Атомарные переменные в Go — это примитивы синхронизации, которые гарантируют выполнение операций **без прерываний**. Они решают проблемы гонок данных (race conditions) и deadlock'ов, обеспечивая безопасный доступ к общим ресурсам из множества горутин.

---

## Что такое атомарные переменные?
**Атомарная операция** — это операция, которая выполняется **за один шаг** и не может быть разделена на более мелкие. Например:
- `atomic.AddInt32(&count, 1)` — атомарный инкремент.
- `atomic.LoadInt32(&count)` — безопасное чтение.

**Зачем нужны?**  
Без атомарности даже простые операции вроде `count++` могут привести к ошибкам:
```go
count++ // Неатомарно! Состоит из:
// 1. Чтение count
// 2. Увеличение на 1
// 3. Запись count
```
## Как работают атомарные операции в Go?
Go предоставляет пакет sync/atomic для атомарных операций.
Основные функции:
```
AddInt32(&var, delta) — атомарное увеличение/уменьшение.
CompareAndSwapInt32(&var, old, new) — проверка и замена значения.
LoadInt32(&var) — безопасное чтение.
StoreInt32(&var, value) — безопасная запись.
```
---
## CAS-операции (Compare-And-Swap)
**CAS** — это механизм, который проверяет, равно ли текущее значение ожидаемому, и если да — заменяет его. 
Например:
```
atomic.CompareAndSwapInt32(&value, old, new)
// 1. Проверяет, что value == old.
// 2. Если да — записывает new.
// 3. Возвращает true/false (успех/неудача).
```
**Зачем нужны?**  
CAS позволяет избежать блокировок (mutex) в некоторых случаях. Например, спин-лок использует CAS для захвата ресурса:
```
func (s *SpinLock) Lock() {
    for !atomic.CompareAndSwapInt32(&s.locked, 0, 1) {
        // Ждем, пока locked станет 0
    }
}
```
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

## Список источников
1. T. Tu, X. Liu, L. Song и Y. Zhang, «Understanding Real-World Concurrency Bugs in Go,» 2019 International Conference on Architectural Support for Programming Languages and Operating Systems (ASPLOS), Providence, RI, USA, апр. 2019, стр. 1–14, doi:10.1145/3297858.3304069.
2. T. Tu, X. Liu, L. Song и Y. Zhang, «Understanding Real-World Concurrency Bugs in Go,» 2023 IEEE International Symposium on Software Reliability Engineering (ISSRE), 2023, стр. 582–592, doi:10.1109/ISSRE62328.2024.00061.
3. evgmih. «Атомики в Go: особенности внутренней реализации». Хабр, 2023. https://habr.com/ru/articles/744822/
4. inetstar. «Go: жарим общие данные. Атомно, быстро и без мьютексов». Хабр, 2024. https://habr.com/ru/company/ruvds/blog/840748/
5. mo0Oonnn. «Погружение в параллелизм в Go». Хабр, 2024. https://habr.com/ru/articles/840750/
6. Zhiyanov, A. «Композиция атомиков в Go». AntonZ.ru, 2024. https://antonz.ru/atomics-composition/
7. Неизвестный автор. «Go FAQ: Какие операции атомарные? Как насчет мьютексов?» Golang блог, 2019. https://golang-blog.blogspot.com/2019/02/go-faq-atomic-ops-mutex.html
8. Go Team. «The Go Memory Model». Официальная документация Go, 2022. https://go.dev/ref/mem
9. Go Team. «sync/atomic: атомарные операции в Go». Официальная документация Go, 2025. https://pkg.go.dev/sync/atomic
10. Caraveo, R. «The Go 1.19 Atomic Wrappers and why to use them». Medium, 2023. https://medium.com/@deckarep/the-go-1-19-atomic-wrappers-and-why-to-use-them-ae14c1177ad8
11. Vincent. «Go: How to Reduce Lock Contention with the Atomic Package». A Journey With Go (Medium), 2020. https://medium.com/a-journey-with-go/go-how-to-reduce-lock-contention-with-the-atomic-package-ba3b2664b549
12. The Quantum Yogi. «The Curious Case of Go’s Memory Model: Simple Language, Subtle Semantics». Medium, 2025. https://medium.com/@kanishksinghpujari/the-curious-case-of-gos-memory-model-simple-language-subtle-semantics-4d3f2029988c
13. Parker, N. «Understanding and Using the sync/atomic Package in Go». Coding Explorations, 2024. https://www.codingexplorations.com/blog/understanding-and-using-the-syncatomic-package-in-go
14. Parker, N. «Understanding Golang's Atomic Package and Mutexes». Coding Explorations, 2023. https://www.codingexplorations.com/blog/understanding-golangs-atomic-package-and-mutexes
15. Dulitha. «Mastering Synchronization Primitives in Go». HackerNoon, 2023. https://hackernoon.com/mastering-synchronization-primitives-in-go
16. Pang. «Is assigning a pointer atomic in Go?» Stack Overflow, 2014. https://stackoverflow.com/questions/21447463/is-assigning-a-pointer-atomic-in-go
17. Drathier. «Is variable assignment atomic in go?» Stack Overflow, 2016. https://stackoverflow.com/questions/33715241/variable-assignment-atomic-in-go
18. api. «Does golang atomic.Load have a acquire semantics?» Stack Overflow, 2019. https://stackoverflow.com/questions/55909553/does-golang-atomic-load-have-an-acquire-semantics
19. Hugh. «Is there a difference in Go between a counter using atomic operations and one using a mutex?» Stack Overflow, 2017. https://stackoverflow.com/questions/47445344/is-there-a-difference-in-go-between-a-counter-using-atomic-operations-and-one-us
20. Yuguang. «Is atomic.LoadUint32 necessary?» Stack Overflow, 2017. https://stackoverflow.com/questions/46556857/is-atomic-loaduint32-necessary
21. Junhee Woo. «Why atomic.Value must not be copied after the first Store?» Stack Overflow, 2021. https://stackoverflow.com/questions/67632767/why-atomic-value-must-not-be-copied-after-the-first-store
