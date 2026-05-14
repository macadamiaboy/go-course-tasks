// Задача 6: Интерфейс для обработки платежа
//
// Ожидаемый вывод:
//   paid 100 by card
//   paid 50 in cash

package main

import "fmt"

// TODO: объяви интерфейс PaymentProcessor с методом Pay(amount int) error

// TODO: объяви структуру CardProcessor (без полей)
// TODO: реализуй Pay(amount int) error для CardProcessor
//       выводи "paid <amount> by card" и возвращай nil

// TODO: объяви структуру CashProcessor (без полей)
// TODO: реализуй Pay(amount int) error для CashProcessor
//       выводи "paid <amount> in cash" и возвращай nil

// TODO: напиши функцию checkout(p PaymentProcessor, amount int)
//       вызывает p.Pay(amount), при ошибке выводит её

func main() {
	// TODO: вызови checkout с CardProcessor{} и amount 100
	// TODO: вызови checkout с CashProcessor{} и amount 50
	fmt.Println("TODO: implement me")
}
