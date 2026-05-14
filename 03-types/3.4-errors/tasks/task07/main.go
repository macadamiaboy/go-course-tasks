// Задача 7: Интеграционная задача — регистрация пользователя
//
// Ожидаемый вывод:
//   registered: Alice (25)
//   error: register: name is empty
//   error: register: age must be >= 18

package main

import "fmt"

// TODO: напиши функцию validateName(name string) error
//       если name == "", верни ошибку "name is empty"

// TODO: напиши функцию validateAge(age int) error
//       если age < 18, верни ошибку "age must be >= 18"

// TODO: напиши функцию register(name string, age int) error
//       вызывает validateName и validateAge
//       оборачивает каждую ошибку через fmt.Errorf: "register: %w"
//       если всё ок — выводи "registered: <name> (<age>)" и возвращай nil

func main() {
	scenarios := []struct {
		name string
		age  int
	}{
		{"Alice", 25},
		{"", 25},
		{"Bob", 15},
	}

	for _, s := range scenarios {
		if err := register(s.name, s.age); err != nil {
			fmt.Println("error:", err)
		}
	}
}
