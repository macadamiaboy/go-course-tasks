// Задача 7: Интеграционная задача — банковский счёт
//
// Ожидаемый вывод:
//   Bob: 250
//   withdraw failed: insufficient funds
//   Bob: 200

package main

import "fmt"

// TODO: объяви структуру BankAccount с полями Owner (string), Balance (int)

// TODO: добавь метод Deposit(amount int) с получателем-указателем

// TODO: добавь метод Withdraw(amount int) bool с получателем-указателем
//       возвращает false если Balance < amount, иначе вычитает и возвращает true

// TODO: добавь метод Info() string с получателем-значением
//       возвращает "<Owner>: <Balance>"

func main() {
	acc := &BankAccount{Owner: "Bob", Balance: 0}
	acc.Deposit(250)
	fmt.Println(acc.Info())

	if !acc.Withdraw(300) {
		fmt.Println("withdraw failed: insufficient funds")
	}

	acc.Withdraw(50)
	fmt.Println(acc.Info())
}
