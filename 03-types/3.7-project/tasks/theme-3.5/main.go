package main

import (
	"fmt"
)

// 3.5.1
// stolen from 3.5 task01
func Echo[T any](v T) T {
	return v
}

func firstTask() {
	fmt.Println("Task 3.5.1:")

	fmt.Println(Echo[int](42))
	fmt.Println(Echo[string]("hello"))
	fmt.Println(Echo[bool](true))
}

// 3.5.2
func First[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	return items[0], true
}

func secondTask() {
	fmt.Println("\nTask 3.5.2:")

	res, ok := First([]int{10, 20, 30})
	fmt.Printf("first: %v, ok: %v\n", res, ok)

	res, ok = First([]int{})
	fmt.Printf("first: %v, ok: %v\n", res, ok)
}

// 3.5.3
type Number interface {
	~int | ~float64
}

func Sum[T Number](items []T) T {
	var total T
	for _, v := range items {
		total += v
	}
	return total
}

func thirdTask() {
	fmt.Println("\nTask 3.5.3:")

	fmt.Println(Sum([]int{1, 2, 3}))
	fmt.Println(Sum([]float64{1.5, 2.5, 6.33}))
}

// 3.5.4
func IndexOf[T comparable](items []T, target T) int {
	for i, val := range items {
		if val == target {
			return i
		}
	}
	return -1
}

func fourthTask() {
	fmt.Println("\nTask 3.5.4:")

	ints := []int{1, 4, 6, 23, 66, 23, 5, 7, 3}
	strings := []string{"asd", "dfgvd", "setfo", "fsvd", "sd"}

	fmt.Println(IndexOf(ints, 5))
	fmt.Println(IndexOf(ints, 23))
	fmt.Println(IndexOf(ints, 9))

	fmt.Println(IndexOf(strings, "sd"))
	fmt.Println(IndexOf(strings, "ssd"))
}

// 3.5.5
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func fifthTask() {
	fmt.Println("\nTask 3.5.5:")

	first := Pair[string, int]{Key: "userID", Value: 194}
	second := Pair[float64, rune]{Key: 9.99, Value: 'p'}
	fmt.Printf("%T: %v\n", first, first)
	fmt.Printf("%T: %v\n", second, second)
}

// 3.5.6
func Values[K comparable, V any](m map[K]V) []V {
	res := []V{}
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

func sixthTask() {
	fmt.Println("\nTask 3.5.6:")

	frequency := make(map[string]int)
	frequency["hey"] = 3
	frequency["honey"] = 9
	frequency["health"] = 19

	res := Values(frequency)

	fmt.Println(res)
}

// 3.5.7
type Store[T any] struct {
	items []T
}

func (s *Store[T]) Add(item T) {
	s.items = append(s.items, item)
}

func (s Store[T]) All() []T { return s.items }

func Contains[T comparable](items []T, target T) bool {
	for _, val := range items {
		if val == target {
			return true
		}
	}
	return false
}

func seventhTask() {
	fmt.Println("\nTask 3.5.7:")

	ints := Store[int]{items: []int{}}
	strings := Store[string]{items: []string{}}

	ints.Add(66)
	ints.Add(67)
	ints.Add(68)
	ints.Add(68)
	fmt.Println("Ints:", ints.All())

	strings.Add("Though")
	strings.Add("Tough")
	strings.Add("Throught")
	strings.Add("Thought")
	strings.Add("Though")
	fmt.Println("Strings:", strings.All())

	fmt.Println("Ints do contain 65:", Contains(ints.All(), 65))
	fmt.Println("Ints do contain 68:", Contains(ints.All(), 68))

	fmt.Println("Strings do contain throughout:", Contains(strings.All(), "Throughout"))
	fmt.Println("Strings do contain thought:", Contains(strings.All(), "Thought"))
}

func main() {
	firstTask()
	secondTask()
	thirdTask()
	fourthTask()
	fifthTask()
	sixthTask()
	seventhTask()
}
