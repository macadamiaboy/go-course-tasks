// Задача 6: Небольшая модель состояния
//
// Ожидаемый вывод:
//   running
//   stopped

package main

import "fmt"

// TODO: объяви структуру Timer с полями Seconds (int), Running (bool)

// TODO: добавь метод Start() с получателем-указателем — ставит Running = true

// TODO: добавь метод Stop() с получателем-указателем — ставит Running = false

// TODO: добавь метод Status() string с получателем-значением
//       возвращает "running" если Running == true, иначе "stopped"

func main() {
	t := &Timer{}
	t.Start()
	fmt.Println(t.Status())
	t.Stop()
	fmt.Println(t.Status())
}
