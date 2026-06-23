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
var ErrNotFound = errors.New("not found")

// TODO: type ValidationError struct { Field, Message string } + метод Error()
type ValidationError struct {
	Field   string
	Message string
}

func (ve ValidationError) Error() string { return fmt.Sprintf("%s: %s", ve.Field, ve.Message) }

// TODO: func repoGetUser(id int) (*User, error)
func repoGetUser(id int) (*User, error) {
	if id == 42 {
		return &User{ID: 42, Name: "John Doe"}, nil
	}
	return nil, fmt.Errorf("get user %d: %w", id, ErrNotFound)
}

// TODO: func serviceGetUser(id int) (*User, error)
func serviceGetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, &ValidationError{Field: "id", Message: "должен быть > 0"}
	}

	user, err := repoGetUser(id)
	if err != nil {
		return nil, fmt.Errorf("repo err: %w", err)
	}
	return user, nil
}

// TODO: func handlerGetUser(id int)
func handlerGetUser(id int) {
	user, err := serviceGetUser(id)
	if err != nil {
		var vErr *ValidationError
		switch {
		case errors.Is(err, ErrNotFound):
			fmt.Println("HTTP 404: not found")
			return
		case errors.As(err, &vErr):
			fmt.Printf("HTTP 400: %s\n", vErr)
			return
		default:
			fmt.Printf("HTTP 500: %s\n", err)
			return
		}
	}
	fmt.Printf("пользователь найден: %s\n", user.Name)
}

func main() {
	// TODO: вызови handlerGetUser(42), handlerGetUser(999), handlerGetUser(-1)
	handlerGetUser(42)
	handlerGetUser(999)
	handlerGetUser(-1)
}
