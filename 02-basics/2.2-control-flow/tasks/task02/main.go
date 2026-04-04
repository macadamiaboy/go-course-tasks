// Задание 2: Подсчёт гласных букв
//
// Напиши функцию countVowels(s string) int,
// которая принимает строку и возвращает количество гласных букв в ней.
// Учитывай и русские, и латинские гласные.
//
// Русские гласные: а, е, ё, и, й, о, у, ы, э, ю, я (и заглавные тоже)
// I believe, the letter "й" is consonant
// Латинские гласные: a, e, i, o, u (и заглавные тоже)
//
// Подсказка: используй for range по строке - он автоматически перебирает
// символы (rune), а не байты. Это важно для кириллицы!
//
// Ожидаемый вывод:
//   "Привет мир" -> гласных: 4
//   You sure?
//   "Hello World" -> гласных: 3
//   "Go" -> гласных: 1
//
// Запусти: go run main.go

package main

import (
	"fmt"
	"slices"
	"strings"
)

// TODO: напиши функцию countVowels(s string) int
// Внутри используй for range и switch для проверки каждого символа

func countVowels(s string) int {
	str := strings.ToLower(s)
	vowels := []rune{'а', 'е', 'ё', 'и', 'й', 'о', 'у', 'ы', 'э', 'ю', 'я', 'a', 'e', 'i', 'o', 'u'}
	counter := 0
	for _, letter := range str {
		if slices.Contains(vowels, letter) {
			counter++
		}
	}
	return counter
}

// at the first made it without the switch expression but reread the task at the skillspace and decided to create another solution
func countVowelsWSwitch(s string) int {
	counter := 0
	for _, letter := range s {
		switch letter {
		case 'а', 'е', 'ё', 'и', 'й', 'о', 'у', 'ы', 'э', 'ю', 'я', 'a', 'e', 'i', 'o', 'u':
			counter++
		case 'А', 'Е', 'Ё', 'И', 'Й', 'О', 'У', 'Ы', 'Э', 'Ю', 'Я', 'A', 'E', 'I', 'O', 'U':
			counter++
		}
	}
	return counter
}

func main() {
	tests := []string{"Привет мир", "Hello World", "Go", "АааыАаАAaAA", "Killerwhale", "", "й"}
	for _, s := range tests {
		fmt.Printf("%q -> гласных: %d\n", s, countVowels(s))
	}

	for _, s := range tests {
		fmt.Printf("%q -> гласных: %d\n", s, countVowelsWSwitch(s))
	}
}
