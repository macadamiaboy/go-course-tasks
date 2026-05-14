// Задание 3: Мемоизация
//
// Напиши функцию memoize(fn func(int) int) func(int) int,
// которая принимает "медленную" функцию и возвращает её кешированную версию.
// Повторный вызов с тем же аргументом должен возвращать результат мгновенно.
//
// Используй замыкание + map[int]int для кеша.
//
// Проверь на функции Фибоначчи:
//   func fib(n int) int {
//       time.Sleep(10 * time.Millisecond) // имитация "медленного" вычисления
//       if n < 2 { return n }
//       return fib(n-1) + fib(n-2)
//   }
//
// Ожидаемый вывод:
//   fib(10) = 55 (вычислено за ~X мс)
//   fib(10) = 55 (из кеша за ~0 мс)
//   fib(20) = 6765 (вычислено за ~X мс)
//
// Запусти: go run main.go

package main

import (
	"fmt"
	"time"
)

// TODO: напиши memoize(fn func(int) int) func(int) int

func fib(n int) int {
	time.Sleep(10 * time.Millisecond)
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func main() {
	// TODO: оберни fib через memoize и замерь время trёх вызовов:
	//   fib(10), повторный fib(10), fib(20)

	_ = fmt.Println
	_ = time.Now
}
