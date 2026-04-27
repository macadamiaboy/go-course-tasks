// Задание 2: Статусы заказа через iota
//
// Тебе нужно:
// 1. Создать константы статусов заказа через iota (пример в README)
// 2. Написать функцию statusName(status int) string,
//    которая возвращает текстовое название статуса
//
// Ожидаемый вывод:
//   Статус 0: Новый
//   Статус 1: В работе
//   Статус 2: Выполнен
//   Статус 3: Отменён
//   Статус 99: Неизвестный статус
//
// Запусти: go run main.go

package main

import "fmt"

// TODO: объяви блок констант с iota для статусов заказа:
// StatusNew, StatusInWork, StatusDone, StatusCancelled
const (
	// TODO: заполни константы здесь
	StatusNew int = iota
	StatusInWork
	StatusDone
	StatusCancelled
)

// TODO: напиши функцию statusName, которая принимает int
// и возвращает строку с названием статуса.
// Для неизвестных значений возвращай "Неизвестный статус".

func statusName(status int) string {
	statuses := [...]string{"StatusNew", "StatusInWork", "StatusDone", "StatusCancelled"}
	if status >= 0 && status < len(statuses) {
		return statuses[status]
	}
	//return "StatusUnknown"
	return "Неизвестный статус"
}

func main() {
	statuses := []int{0, 1, 2, 3, 4, -1}
	for _, s := range statuses {
		fmt.Printf("Статус %d: %s\n", s, statusName(s))
	}
}
