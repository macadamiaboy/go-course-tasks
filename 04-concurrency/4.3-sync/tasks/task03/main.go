// Задание 3: Atomic vs Mutex - два подхода к счётчику
//
// Реализуй счётчик двумя способами и сравни их.
//
// Способ 1 - через sync.Mutex:
//   type MutexCounter struct { mu sync.Mutex; value int }
//   func (c *MutexCounter) Increment()
//   func (c *MutexCounter) Value() int
//
// Способ 2 - через sync/atomic (Go 1.19+):
//   type AtomicCounter struct { value atomic.Int64 }
//   func (c *AtomicCounter) Increment()
//   func (c *AtomicCounter) Value() int64
//
// В main() запусти 1000 горутин для каждого счётчика.
// Оба должны вернуть 1000.
//
// Проверь: go run -race main.go - предупреждений быть не должно.
//
// Ожидаемый вывод:
//   MutexCounter:  1000
//   AtomicCounter: 1000
//
// Запусти: go run -race main.go

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// TODO: реализуй MutexCounter (sync.Mutex)

// TODO: реализуй AtomicCounter (atomic.Int64)

func main() {
	// TODO: запусти 1000 горутин для каждого счётчика,
	// дождись завершения через WaitGroup и напечатай значения.

	_ = fmt.Println
	_ = sync.WaitGroup{}
	_ = atomic.Int64{}
}
