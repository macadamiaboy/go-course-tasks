// Задача 2: Нулевые значения
//
// Ожидаемый вывод:
//   { 0 false}
//   Host:
//   Port: 0
//   Debug: false

package main

import "fmt"

// TODO: объяви структуру AppConfig с полями Host (string), Port (int), Debug (bool)

func main() {
	// TODO: объяви var cfg AppConfig без инициализации

	fmt.Println(cfg)
	fmt.Println("Host:", cfg.Host)
	fmt.Println("Port:", cfg.Port)
	fmt.Println("Debug:", cfg.Debug)
}
