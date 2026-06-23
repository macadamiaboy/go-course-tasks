package task01

import (
	"math"
	"testing"
)

// TODO: напиши табличные тесты для функции ParsePrice.
//
// Создай срез структур с полями:
//   name        string   - название теста
//   input       string   - входная строка
//   expected    float64  - ожидаемый результат
//   expectError bool     - ожидаем ли ошибку
//
// Обязательно проверь эти случаи:
//   Корректные: "1500", "1 500", "1500.50", "1500,50", "1 500,50 руб.", "₽1 500"
//   Граничные:  "0", "0,01"
//   Некорректные: "", "abc", "руб.", "---"
//
// Подсказка: запусти go test -v ./... чтобы увидеть подробный вывод

func TestParsePrice(t *testing.T) {
	cases := []struct {
		name        string
		input       string
		expected    float64
		expectError bool
	}{
		// TODO: заполни таблицу тест-кейсов
		// корректные
		{"корректное целое", "1500", 1500.0, false},
		{"корректное целое с пробелом", "1 500", 1500.0, false},
		{"корректное дробное", "1500.50", 1500.50, false},
		{"корректное дробное с запятой", "1500,50", 1500.50, false},
		{"корректное дробное с запятой и пробелом", "1 500,50", 1500.50, false},
		{"корректное с пробелом и валютой", "₽1 500", 1500.0, false},
		// граничные
		{"граничное ноль", "0", 0, false},
		{"граничное с запятой", "0,01", 0.01, false},
		// некорректные
		{"некорректное отсутствующее", "", 0, true},
		{"некорректное буквы", "abc", 0, true},
		{"некорректное буквы русские с точкой", "руб.", 0, true},
		{"некорректное символы", "---", 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// TODO: вызови ParsePrice(tc.input)
			res, err := ParsePrice(tc.input)
			// TODO: проверь что err != nil если expectError == true
			if err != nil && !tc.expectError {
				t.Errorf("didn't expect the error but got: %s", err.Error())
			}
			// TODO: проверь что результат == tc.expected (используй сравнение с погрешностью 0.001)
			epsilon := 0.01
			if math.Abs(res-tc.expected) > epsilon {
				t.Errorf("the result didn't match the expectations. Expected: %f.2, got: %f.2", tc.expected, res)
			}
		})
	}
}

// FuzzParsePrice проверяет что функция не паникует ни при каких входных данных.
// Запуск: go test -fuzz=FuzzParsePrice -fuzztime=10s
func FuzzParsePrice(f *testing.F) {
	// TODO: добавь начальные случаи через f.Add(...)
	// f.Add("1500")
	// f.Add("")
	f.Add("1500")
	f.Add("1500,50 руб.")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		// TODO: вызови ParsePrice(input)
		// Функция не должна паниковать - только возвращать ошибку
		// Не нужно проверять результат - главное что нет panic
		ParsePrice(input)
	})
}
