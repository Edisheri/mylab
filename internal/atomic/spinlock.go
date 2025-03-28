package atomic

import (
	"runtime"
	"sync/atomic"
)

type SpinLock struct {
	locked int32 // 0 — свободен, 1 — занят
}

func (s *SpinLock) Lock() {
	// Пытаемся атомарно изменить locked с 0 на 1
	for !atomic.CompareAndSwapInt32(&s.locked, 0, 1) {
		runtime.Gosched() // Передаём управление другим горутинам
	}
}

func (s *SpinLock) Unlock() {
	atomic.StoreInt32(&s.locked, 0) // Освобождаем блокировку
}
