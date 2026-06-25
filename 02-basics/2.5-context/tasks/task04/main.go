// Задание 4: Контекст через три слоя
//
// Реализуй цепочку handler -> service -> repository, где каждая функция
// принимает ctx context.Context первым параметром.
//
// 1. Объяви свой тип ключа и константу:
//      type ctxKey string
//      const requestIDKey ctxKey = "request-id"
//
// 2. handler(ctx): кладёт в контекст requestID="req-99" через context.WithValue
//    и вызывает service(ctx).
//
// 3. service(ctx): вызывает repository(ctx). Может добавить WithTimeout.
//
// 4. repository(ctx):
//      - достаёт requestID из контекста
//      - печатает "[req-99] получаем данные из БД..."
//      - имитирует работу через time.Sleep
//      - проверяет ctx.Done(): если контекст отменён - возвращает ctx.Err()
//      - иначе печатает "[req-99] готово" и возвращает nil
//
// Ожидаемый вывод:
//   [req-99] получаем данные из БД...
//   [req-99] готово
//
// Запусти: go run main.go

package main

import (
	"context"
	"fmt"
	"time"
)

// TODO: объяви тип ключа и константу requestIDKey
type contextKey string

const requestIDKey contextKey = "request-id"

// TODO: handler(ctx context.Context) error
func handler(ctx context.Context) {
	requestID := "req-99"
	wrappedCtx := context.WithValue(ctx, requestIDKey, requestID)
	if err := service(wrappedCtx); err != nil {
		fmt.Println("ошибка в handler:", err)
	}
}

// TODO: service(ctx context.Context) error
func service(ctx context.Context) error {
	err := repository(ctx)
	if err != nil {
		return fmt.Errorf("service: %w", err)
	}
	return nil
}

// TODO: repository(ctx context.Context) error
func repository(ctx context.Context) error {
	requestID := ctx.Value(requestIDKey)
	fmt.Printf("[%v] получаем данные из БД...\n", requestID)

	time.Sleep(50 * time.Millisecond)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		fmt.Printf("[%v] готово\n", requestID)
		return nil
	}
}

func main() {
	// TODO: вызови handler(context.Background()) и обработай ошибку
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handler(ctx)
}
