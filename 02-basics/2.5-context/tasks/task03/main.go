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
//
//	case <-ctx.Done(): печатай "воркер N остановлен: <ctx.Err()>" и return
//	case <-ticker.C:  печатай "воркер N работает..."
func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			fmt.Printf("воркер %d работает...\n", id)
		case <-ctx.Done():
			fmt.Printf("воркер %d остановлен: %v\n", id, ctx.Err())
			return
		}
	}
}

func main() {
	// TODO: создай контекст с отменой
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// TODO: запусти двух воркеров в горутинах, дождись их через sync.WaitGroup
	// TODO: через 2 секунды вызови cancel()
	wg.Add(1)
	go worker(ctx, 1, &wg)
	wg.Add(1)
	go worker(ctx, 2, &wg)

	time.Sleep(2 * time.Second)
	cancel()
	wg.Wait()
}
