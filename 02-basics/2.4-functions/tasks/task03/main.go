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
// по-моему, не эффективно будет, так как не кэшируются промежуточные шаги
// если у нас в кэше будут значения до 5, а затем мы вызовем fn(10), то запишем в кэш только 10
// с 5 по 10 останутся пробелы
func memoize(fn func(int) int) func(int) int {
	cache := map[int]int{}
	return func(i int) int {
		if i, ok := cache[i]; ok {
			return i
		}
		result := fn(i)
		cache[i] = result
		return result
	}
}

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
	/*
		mFib := memoize(fib)

		start := time.Now()
		mFib(10)
		first := time.Since(start).Milliseconds()

		start = time.Now()
		mFib(10)
		second := time.Since(start).Milliseconds()

		start = time.Now()
		mFib(20)
		third := time.Since(start).Milliseconds()

		fmt.Printf("First call: %dms, second call: %dms, third call: %dms\n", first, second, third)
	*/

	var modernized func(int) int
	modernized = memoize(func(n int) int {
		time.Sleep(10 * time.Millisecond)
		if n < 2 {
			return n
		}
		return modernized(n-1) + modernized(n-2)
	})

	start := time.Now()
	modernized(10)
	first := time.Since(start).Milliseconds()

	start = time.Now()
	modernized(10)
	second := time.Since(start).Milliseconds()

	start = time.Now()
	modernized(20)
	third := time.Since(start).Milliseconds()

	fmt.Printf("First call: %dms, second call: %dms, third call: %dms\n", first, second, third)
}
