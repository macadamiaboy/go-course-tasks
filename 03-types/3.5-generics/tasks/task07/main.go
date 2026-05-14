// Задача 7: Интеграционная задача — обобщенный Store
//
// Ожидаемый вывод:
//   strings: [go rust python]
//   ints: [1 2 3]
//   contains "go": true
//   contains "java": false
//   contains 2: true
//   contains 5: false

package main

import "fmt"

// TODO: объяви структуру Store[T any] с полем items []T
// TODO: добавь метод Add(item T) — добавляет элемент в items
// TODO: добавь метод All() []T — возвращает items

// TODO: напиши функцию Contains[T comparable](items []T, target T) bool
//       возвращает true если target есть в items

func main() {
	// TODO: создай Store[string], добавь "go", "rust", "python"
	//       выведи "strings:", ss.All()

	// TODO: создай Store[int], добавь 1, 2, 3
	//       выведи "ints:", si.All()

	// TODO: выведи результаты Contains:
	//       Contains(ss.All(), "go")  → "contains \"go\": true"
	//       Contains(ss.All(), "java") → "contains \"java\": false"
	//       Contains(si.All(), 2)     → "contains 2: true"
	//       Contains(si.All(), 5)     → "contains 5: false"
	fmt.Println("TODO: implement me")
}
