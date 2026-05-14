// Задание 3: Остановить воркер через контекст
//
// Напиши функцию worker(ctx context.Context, id int), которая:
//   - каждые 500ms печатает "воркер N работает..."
//   - завершается когда ctx.Done() закрывается и печатает
//     "воркер N остановлен: <ctx.Err()>"
//
// В main() запусти двух воркеров и останови их через 2 секунды
// через context.WithCancel.
//
// Ожидаемый вывод (порядок между воркерами может отличаться):
//   воркер 1 работает...
//   воркер 2 работает...
//   воркер 1 работает...
//   воркер 2 работает...
//   воркер 1 остановлен: context canceled
//   воркер 2 остановлен: context canceled
//
// Запусти: go run main.go

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TODO: напиши функцию worker(ctx context.Context, id int)
// Внутри select:
//   case <-ctx.Done(): печатай "воркер N остановлен: <ctx.Err()>" и return
//   case <-ticker.C:  печатай "воркер N работает..."

func main() {
	// TODO: создай контекст с отменой
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// TODO: запусти двух воркеров в горутинах, дождись их через sync.WaitGroup
	// TODO: через 2 секунды вызови cancel()

	_ = context.Background()
	_ = fmt.Println
	_ = sync.WaitGroup{}
	_ = time.Second
}
