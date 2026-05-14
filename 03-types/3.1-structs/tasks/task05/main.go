// Задача 5: Композиция с Audit
//
// Ожидаемый вывод:
//   Title: Go Basics
//   CreatedAt: 2024-01-01
//   UpdatedAt: 2024-06-01

package main

import "fmt"

// TODO: объяви структуру Audit с полями CreatedAt (string), UpdatedAt (string)

// TODO: объяви структуру Article с полем Title (string) и встроенным полем Audit (embedding)

func main() {
	// TODO: создай значение Article:
	// Title: "Go Basics"
	// Audit: {CreatedAt: "2024-01-01", UpdatedAt: "2024-06-01"}

	fmt.Println("Title:", a.Title)
	fmt.Println("CreatedAt:", a.CreatedAt)
	fmt.Println("UpdatedAt:", a.UpdatedAt)
}
