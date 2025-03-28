package atomic

import (
	"sync/atomic"
)

// SafeCounter безопасно увеличивается/уменьшается из нескольких горутин
type SafeCounter struct {
	value int32
}

// Increment увеличивает счетчик на 1
func (c *SafeCounter) Increment() {
	atomic.AddInt32(&c.value, 1)
}

// Decrement уменьшает счетчик на 1
func (c *SafeCounter) Decrement() {
	atomic.AddInt32(&c.value, -1)
}

// Value возвращает текущее значение счетчика
func (c *SafeCounter) Value() int32 {
	return atomic.LoadInt32(&c.value)
}
