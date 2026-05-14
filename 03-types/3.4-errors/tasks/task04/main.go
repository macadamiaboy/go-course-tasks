// Задача 4: Оборачивание ошибки
//
// Ожидаемый вывод:
//   load data: read file: file not found

package main

import "fmt"

// TODO: напиши функцию readFileMock() error
//       возвращай errors.New("file not found")

// TODO: напиши функцию loadData() error
//       вызывает readFileMock() и оборачивает ошибку через fmt.Errorf с %w:
//       "read file: <err>"
//       сама тоже оборачивает: "load data: <err>"

func main() {
	err := loadData()
	fmt.Println(err)
}
