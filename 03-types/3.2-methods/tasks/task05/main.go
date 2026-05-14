// Задача 5: Безопасная работа с nil
//
// Ожидаемый вывод:
//   user is nil
//   Alice

package main

import "fmt"

// TODO: объяви структуру User с полем Name (string)

// TODO: напиши функцию printUserName(u *User)
//       если u == nil — выведи "user is nil"
//       иначе — выведи u.Name

func main() {
	printUserName(nil)
	printUserName(&User{Name: "Alice"})
}
