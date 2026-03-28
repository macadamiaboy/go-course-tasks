// Задание 2: Инвертировать словарь
//
// Напиши функцию invertMap, которая:
//   - принимает map[string]int (название -> число)
//   - возвращает map[int]string (число -> название)
//
// Потом выведи результат: отсортируй ключи по возрастанию и выведи каждую пару.
//
// Используй пакет slices для сортировки ключей.
//
// Ожидаемый вывод:
//   1 -> яблоко
//   2 -> банан
//   3 -> апельсин
//
// Запусти: go run main.go

package main

import (
	"fmt"
	"maps"
	"slices"
)

// TODO: напиши функцию invertMap(m map[string]int) map[int]string

func invertMap(m map[string]int) map[int]string {
	resultMap := make(map[int]string)

	for key, value := range m {
		resultMap[value] = key
	}
	return resultMap
}

func main() {
	fruits := map[string]int{
		"яблоко":   1,
		"банан":    2,
		"апельсин": 3,
		"груша":    10,
		"киви":     5,
		"шоколад":  2,
	}

	// TODO: вызови invertMap и сохрани результат
	// inverted := invertMap(fruits)

	// TODO: собери ключи из inverted в срез, отсортируй их
	// и выведи каждую пару в формате: "1 -> яблоко"

	inverted := invertMap(fruits)
	resultSlice := slices.Collect(maps.Keys(inverted))
	slices.Sort(resultSlice)

	for _, element := range resultSlice {
		fmt.Printf("%d -> %s\n", element, inverted[element])
	}

	// Забавно, в результате случайным образом могут по ключу 2 вернуться и банан, и шоколад, так как в методе invertMap
	// пары значений берутся рандомно. И в зависимости от того, что придет вторым, то и останется в инвертированной мапе
}
