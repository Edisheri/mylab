package atomic

import "sync/atomic"

type AtomicFlag struct {
	flag int32 // 0 — false, 1 — true
}

func (f *AtomicFlag) Set() {
	atomic.StoreInt32(&f.flag, 1) // Атомарно устанавливает флаг в 1
}

func (f *AtomicFlag) Unset() {
	atomic.StoreInt32(&f.flag, 0) // Атомарно сбрасывает флаг в 0
}

func (f *AtomicFlag) IsSet() bool {
	return atomic.LoadInt32(&f.flag) == 1 // Проверяет, установлен ли флаг
}
