package main

import (
	"errors"
	"fmt"
	"strconv"
)

// 3.4.1
// stolen from 3.4 task01
func safeDivide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("dividion by zero")
	}
	return a / b, nil
}

func divideCall(a, b int) {
	res, err := safeDivide(a, b)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("result:", res)
}

func firstTask() {
	fmt.Println("Task 3.4.1:")

	divideCall(10, 2)
	divideCall(10, 0)
}

// 3.4.2
func createUser(name string) error {
	if name == "" {
		return errors.New("name is empty")
	}
	return nil
}

func createUserCaller(name string) {
	err := createUser(name)
	if err != nil {
		fmt.Println("Error occurred:", err)
	} else {
		fmt.Println("Everything's ok")
	}
}

func secondTask() {
	fmt.Println("\nTask 3.4.2:")

	createUserCaller("Alex")
	createUserCaller("")
}

// 3.4.3
var ErrOutOfStock = errors.New("out of stock")

func buyItem(count int) error {
	if count == 0 {
		return ErrOutOfStock
	} else if count < 0 {
		return errors.New("Negative argument")
	}
	return nil
}

func buyItemCaller(count int) {
	err := buyItem(count)
	if err != nil {
		if errors.Is(err, ErrOutOfStock) {
			fmt.Println("Sentinel:", err)
			return
		}
		fmt.Println("Unknown error occurred:", err)
		return
	}
	fmt.Println("Everything's ok")
}

func thirdTask() {
	fmt.Println("\nTask 3.4.3:")

	buyItemCaller(99)
	buyItemCaller(0)
	buyItemCaller(-1)
}

// 3.4.4
func readFileMock() error { return errors.New("base error") }

func loadData() error {
	err := readFileMock()
	return fmt.Errorf("Wrapped error: %w", err)
}

func fourthTask() {
	fmt.Println("\nTask 3.4.4:")

	fmt.Println(loadData())
}

// 3.4.5
type InputError struct {
	Field  string
	Reason string
}

func (ie InputError) Error() string { return fmt.Sprintf("%s: %s", ie.Field, ie.Reason) }

func validateEmail(email string) error {
	if email == "" {
		return InputError{Field: "email", Reason: "is empty"}
	}
	return nil
}

func fifthTask() {
	fmt.Println("\nTask 3.4.5:")

	var err InputError
	errors.As(validateEmail(""), &err)
	fmt.Printf("%T: %s\n", err, err)
}

// 3.4.6
var IDs []string = []string{"123", "43156", "23634", "", "32084", ""}

func parseID(s string) (int, error) {
	if s == "" {
		return 0, errors.New("the id is empty")
	}

	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func sixthTask() {
	fmt.Println("\nTask 3.4.6:")

	fmt.Println(IDs)
	for _, id := range IDs {
		_, err := parseID(id)
		if err != nil {
			fmt.Print("skip ")
			continue
		}
		fmt.Print("ok ")
	}
	fmt.Println()
}

// 3.4.7
func validateName(name string) error {
	if name == "" {
		return errors.New("Name should not be empty")
	}
	return nil
}

func validateAge(age int) error {
	if age <= 0 {
		return errors.New("Age cannot be negative or zero")
	}
	return nil
}

func register(name string, age int) error {
	if err := validateName(name); err != nil {
		return fmt.Errorf("Failed to validate the name, %w", err)
	}

	if err := validateAge(age); err != nil {
		return fmt.Errorf("Failed to validate the age, %w", err)
	}

	return nil
}

func seventhTask() {
	fmt.Println("\nTask 3.4.7:")

	// Both illegal args. It will ignore the second
	fmt.Println(register("", -1))
	// Illegal name
	fmt.Println(register("", 19))
	// Illegal age
	fmt.Println(register("Kevin", -19))
	// Legal args
	fmt.Println(register("Arvid", 19))
}

func main() {
	firstTask()
	secondTask()
	thirdTask()
	fourthTask()
	fifthTask()
	sixthTask()
	seventhTask()
}
