// Задача 4: Поиск индекса
//
// Ожидаемый вывод:
//   index: 2
//   index: -1

package main

import "fmt"

// TODO: напиши функцию IndexOf[T comparable](items []T, target T) int
//
//	возвращает индекс target в items, или -1 если не найден
func IndexOf[T comparable](items []T, target T) int {
	for i, element := range items {
		if element == target {
			return i
		}
	}
	return -1
}

func main() {
	// TODO: вызови IndexOf([]string{"a", "b", "c"}, "c") и выведи "index: <результат>"
	// TODO: вызови IndexOf([]string{"a", "b", "c"}, "z") и выведи "index: <результат>"
	fmt.Printf("index: %d\n", IndexOf([]string{"a", "b", "c"}, "c"))
	fmt.Printf("index: %d\n", IndexOf([]string{"a", "b", "c"}, "z"))
}
