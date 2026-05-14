// Задание 2: Цепочка ошибок (repository → service → handler)
//
// Имитируем трёхслойное приложение.
//
// 1. Sentinel-ошибка уровня repository:
//      var ErrNotFound = errors.New("not found")
//
// 2. Типизированная ошибка уровня service:
//      type ValidationError struct { Field, Message string }
//      (метод Error() string)
//
// 3. Слои:
//      - repository.GetUser(id int) (*User, error)
//          id == 42 → возвращает пользователя
//          иначе    → fmt.Errorf("get user %d: %w", id, ErrNotFound)
//
//      - service.GetUser(id int) (*User, error)
//          id <= 0  → &ValidationError{Field: "id", Message: "должен быть > 0"}
//          иначе    → вызывает repository и оборачивает ошибку через %w
//
//      - handler.GetUser(id int)
//          вызывает service
//          если errors.Is(err, ErrNotFound)        → выводит "HTTP 404: not found"
//          если errors.As(err, &ValidationError{}) → выводит "HTTP 400: <поле> <сообщение>"
//          иначе                                   → "HTTP 500: <err>"
//
// В main() вызови handler.GetUser с тремя значениями: 42 (ok), 999 (не найден), -1 (валидация).
//
// Ожидаемый вывод:
//   id=42: пользователь найден: Аня
//   id=999: HTTP 404: not found
//   id=-1: HTTP 400: поле "id" - должен быть > 0
//
// Запусти: go run main.go

package main

import (
	"errors"
	"fmt"
)

type User struct {
	ID   int
	Name string
}

// TODO: var ErrNotFound = errors.New("not found")

// TODO: type ValidationError struct { Field, Message string } + метод Error()

// TODO: func repoGetUser(id int) (*User, error)

// TODO: func serviceGetUser(id int) (*User, error)

// TODO: func handlerGetUser(id int)

func main() {
	// TODO: вызови handlerGetUser(42), handlerGetUser(999), handlerGetUser(-1)

	_ = errors.Is
	_ = errors.As
	_ = fmt.Println
}
