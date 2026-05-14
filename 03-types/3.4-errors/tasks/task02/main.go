// Задача 2: Проверка пустого имени
//
// Ожидаемый вывод:
//   user created: Alice
//   error: name is empty

package main

import (
	"fmt"
)

// TODO: напиши функцию createUser(name string) error
//       если name == "", верни errors.New("name is empty")
//       иначе выведи "user created: <name>" и верни nil
//       (не забудь импортировать "errors")

func main() {
	if err := createUser("Alice"); err != nil {
		fmt.Println("error:", err)
	}
	if err := createUser(""); err != nil {
		fmt.Println("error:", err)
	}
}
