// Задача 7: Интеграционная мини-модель
//
// Ожидаемый вывод:
//   Enrollment E-001: Alice enrolled in Go Basics [active]

package main

import "fmt"

// TODO: объяви структуру ContactInfo с полями Phone (string), Email (string)

// TODO: объяви структуру Course с полями ID (string), Title (string)

// TODO: объяви структуру Student с полями Name (string) и встроенным ContactInfo (embedding)

// TODO: объяви структуру CourseEnrollment с полями:
//       EnrollmentID (string), Status (string)
//       Student Student, Course Course

func main() {
	// TODO: создай значение CourseEnrollment с осмысленными данными

	fmt.Printf("Enrollment %s: %s enrolled in %s [%s]\n",
		e.EnrollmentID, e.Student.Name, e.Course.Title, e.Status)
}
