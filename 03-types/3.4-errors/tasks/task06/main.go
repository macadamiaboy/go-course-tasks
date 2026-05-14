// Задача 6: Обработка ошибок в цикле
//
// Ожидаемый вывод:
//   ok: 1
//   skip
//   ok: 2
//   skip
//   ok: 3

package main

import (
	"fmt"
	"strconv"
)

// TODO: напиши функцию parseID(s string) (int, error)
//       если s == "", верни 0 и errors.New("empty id")
//       иначе используй strconv.Atoi(s) и верни результат

func main() {
	ids := []string{"1", "", "2", "", "3"}

	for _, id := range ids {
		// TODO: вызови parseID(id)
		//       если ошибка — выведи "skip" и продолжи
		//       если успех — выведи "ok: <значение>"
		_ = id
		_ = strconv.Atoi // подсказка: используй в parseID
		fmt.Println("TODO: implement me")
	}
}
