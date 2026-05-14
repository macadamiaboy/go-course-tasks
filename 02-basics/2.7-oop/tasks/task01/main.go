// Задание 1: Стек
//
// Реализуй структуру Stack для int с методами:
//   - Push(val int)           - добавить элемент на вершину
//   - Pop() (int, error)      - забрать верхний элемент; ошибка, если стек пуст
//   - Peek() (int, error)     - посмотреть верхний элемент без удаления
//   - Len() int               - количество элементов
//   - IsEmpty() bool          - true, если стек пуст
//
// Методы должны быть определены на *Stack (указателе), чтобы Push/Pop
// изменяли внутреннее состояние.
//
// Заведи переменную ErrEmptyStack через errors.New — её будут возвращать Pop/Peek.
//
// В main() продемонстрируй работу: запушь 3 числа, вызови Peek, сделай Pop
// до опустошения и один Pop на пустом стеке — он должен вернуть ошибку.
//
// Ожидаемый вывод:
//   Len=3 Top=30
//   Pop: 30
//   Pop: 20
//   Pop: 10
//   Pop пустого: stack is empty
//
// Запусти: go run main.go

package main

import (
	"errors"
	"fmt"
)

// TODO: var ErrEmptyStack = errors.New("stack is empty")

// TODO: type Stack struct { items []int }

// TODO: методы Push, Pop, Peek, Len, IsEmpty на *Stack

func main() {
	// TODO: s := &Stack{}
	// s.Push(10); s.Push(20); s.Push(30)
	// ...

	_ = errors.New
	_ = fmt.Println
}
