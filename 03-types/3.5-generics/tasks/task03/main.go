// Задача 3: Сумма чисел
//
// Ожидаемый вывод:
//   sum int: 15
//   sum float: 7.5

package main

import "fmt"

// TODO: объяви ограничение Number для ~int | ~float64
type Number interface {
	~int | ~float64
}

// TODO: напиши функцию Sum[T Number](items []T) T
//
//	возвращает сумму всех элементов среза
func Sum[T Number](items []T) T {
	var res T
	for _, element := range items {
		res += element
	}
	return res
}

func main() {
	// TODO: вызови Sum для []int{1, 2, 3, 4, 5} и выведи "sum int: <результат>"
	// TODO: вызови Sum для []float64{1.5, 2.0, 4.0} и выведи "sum float: <результат>"
	ints := []int{1, 2, 3, 4, 5}
	floats := []float64{1.5, 2.0, 4.0}

	fmt.Printf("sum int: %d\n", Sum(ints))
	fmt.Printf("sum float: %.1f\n", Sum(floats))
}
