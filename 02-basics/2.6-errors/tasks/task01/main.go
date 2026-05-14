// Задание 1: Типизированная ошибка валидации
//
// Реализуй тип ValidationError с полями Field и Message.
// Он должен удовлетворять интерфейсу error (метод Error() string).
//
// Напиши функцию validateUser(name, email string) error:
//   - если name == ""           -> *ValidationError{Field: "name",  Message: "обязательное поле"}
//   - если в email нет символа '@' -> *ValidationError{Field: "email", Message: "неверный формат (нет @)"}
//   - иначе -> nil
//
// В main() вызови validateUser с заведомо невалидными данными,
// поймай ошибку через errors.As и выведи поля.
//
// Ожидаемый вывод:
//   ошибка валидации: поле "email" - неверный формат (нет @)
//
// Запусти: go run main.go

package main

import (
	"errors"
	"fmt"
)

// TODO: определи тип ValidationError и метод Error() string

// TODO: напиши функцию validateUser(name, email string) error

func main() {
	// TODO: вызови validateUser с невалидным email, например validateUser("Аня", "anya.mail")
	//
	// var vErr *ValidationError
	// if errors.As(err, &vErr) {
	//     fmt.Printf(...)
	// }

	_ = errors.As // убери когда начнёшь использовать errors
	_ = fmt.Println
}
