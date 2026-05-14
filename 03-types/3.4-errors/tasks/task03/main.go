// Задача: Sentinel error
//
// Ожидаемый вывод:
//   buy ok
//   out of stock: show restock page

package main

import "fmt"

// TODO: объяви sentinel error:
//       var ErrOutOfStock = errors.New("out of stock")
//       (не забудь импортировать "errors")

// TODO: напиши функцию buyItem(count int) error
//       если count == 0, верни ErrOutOfStock
//       иначе верни nil

func main() {
	// TODO: вызови buyItem(5)
	//       если ошибки нет — выведи "buy ok"

	// TODO: вызови buyItem(0)
	//       если errors.Is(err, ErrOutOfStock) — выведи "out of stock: show restock page"
	fmt.Println("TODO: implement me")
}
