// Задача: Деление с проверкой
//
// Ожидаемый вывод:
//   result: 5
//   error: division by zero

package main

import (
	"errors"
	"fmt"
)

// TODO: напиши функцию safeDivide(a, b int) (int, error)
//
//	если b == 0, верни 0 и ошибку "division by zero"
//	иначе верни a / b и nil
func safeDivide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// tried to do it not hardcoded and as idiomatic as possible
func divideCall(a, b int) {
	res, err := safeDivide(a, b)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("result:", res)
}

func main() {
	// TODO: вызови safeDivide(10, 2) и выведи "result: <значение>"
	divideCall(10, 2)
	// TODO: вызови safeDivide(10, 0) и выведи "error: <ошибка>"
	divideCall(10, 0)
}
