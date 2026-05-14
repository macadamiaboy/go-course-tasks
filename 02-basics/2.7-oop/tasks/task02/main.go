// Задание 2: Фигуры через интерфейс
//
// Создай интерфейс Shape с методами:
//   - Area() float64
//   - Perimeter() float64
//
// Реализуй его для двух типов:
//   - Circle    { Radius float64 }       -> πr², 2πr
//   - Rectangle { Width, Height float64 }-> W*H, 2*(W+H)
//
// Используй math.Pi для π.
//
// Напиши функцию printShapeInfo(name string, s Shape), которая печатает
// площадь и периметр в формате:
//   <name>: площадь=<area>, периметр=<perimeter>
// с двумя знаками после запятой (формат %.2f).
//
// В main() создай Circle{Radius: 5} и Rectangle{Width: 4, Height: 6}
// и вызови printShapeInfo для каждой.
//
// Ожидаемый вывод:
//   Круг: площадь=78.54, периметр=31.42
//   Прямоугольник: площадь=24.00, периметр=20.00
//
// Запусти: go run main.go

package main

import (
	"fmt"
	"math"
)

// TODO: type Shape interface { Area() float64; Perimeter() float64 }

// TODO: type Circle struct { Radius float64 } + методы Area, Perimeter

// TODO: type Rectangle struct { Width, Height float64 } + методы Area, Perimeter

// TODO: func printShapeInfo(name string, s Shape)

func main() {
	// TODO: printShapeInfo("Круг", Circle{Radius: 5})
	// TODO: printShapeInfo("Прямоугольник", Rectangle{Width: 4, Height: 6})

	_ = fmt.Println
	_ = math.Pi
}
