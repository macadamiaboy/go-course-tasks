package main

import (
	"fmt"
	"math"
)

// 3.3.1
type Greeter interface {
	Greet() string
}

type User struct {
	Name string
}

type Guest struct{}

func (u User) Greet() string { return fmt.Sprintf("Hello, %s", u.Name) }

func (g Guest) Greet() string { return "Welcome, stranger" }

func printGreeting(g Greeter) { fmt.Println(g.Greet()) }

func firstTask() {
	fmt.Println("Task 3.3.1:")
	printGreeting(User{Name: "Anna"})
	printGreeting(Guest{})
}

// 3.3.2
type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 { return r.Width * r.Height }

func (c Circle) Area() float64 { return math.Pi * math.Pow(c.Radius, 2) }

func secondTask() {
	fmt.Println("\nTask 3.3.2:")

	shapes := []Shape{Rectangle{Width: 10, Height: 5}, Circle{Radius: 3}}
	for _, shape := range shapes {
		area := shape.Area()
		switch shape.(type) {
		case Rectangle:
			fmt.Printf("Rectangle area: %.2f\n", area)
		case Circle:
			fmt.Printf("Circle area: %.2f\n", area)
		}
	}
}

// 3.3.3
type Runner interface {
	Run() string
}

type Athlete struct {
	Sport string
	Name  string
}

func (a Athlete) Run() string { return fmt.Sprintf("%s is running", a.Name) }

func thirdTask() {
	fmt.Println("\nTask 3.3.3:")

	var _ Runner = Athlete{}
	fmt.Println("Everything's ok")
}

// 3.3.4
type Starter interface {
	Start() string
}

type Stopper interface {
	Stop() string
}

type Machine interface {
	Starter
	Stopper
}

type Engine struct {
	Name      string
	Fuel      string
	Type      string
	Cylinders int
}

func (e Engine) Start() string { return fmt.Sprintf("%s go ratata", e.Name) }

func (e Engine) Stop() string { return fmt.Sprintf("%s is stopped", e.Name) }

func runCycle(m Machine) {
	fmt.Println(m.Start())
	fmt.Println(m.Stop())
}

func fourthTask() {
	fmt.Println("\nTask 3.3.4:")

	engine := Engine{Name: "s85b50", Fuel: "petrol", Type: "V-type", Cylinders: 10}
	runCycle(engine)
}

// 3.3.5
func printStringLength(x any) {
	str, ok := x.(string)
	if !ok {
		fmt.Println("not a string")
		return
	}
	fmt.Println("The length of the string is", len([]rune(str)))
}

func fifthTask() {
	fmt.Println("\nTask 3.3.5:")

	fmt.Println("Trying on an int")
	printStringLength(15)

	fmt.Println("Trying on a string")
	printStringLength("throughout")

	fmt.Println("Trying on an int using conversion")
	integer := 15
	printStringLength(string(integer))
	fmt.Println("Go treats int like a rune and creates a string of one rune not a string of digits")
}

// 3.3.6
type PaymentProcessor interface {
	Pay(int) error
}

type CardProcessor struct{}

type CashProcessor struct{}

func (CardProcessor) Pay(amount int) error {
	fmt.Printf("%d dollars are payed by card\n", amount)
	return nil
}

func (CashProcessor) Pay(amount int) error {
	fmt.Printf("%d dollars are payed by cash\n", amount)
	return nil
}

func checkout(p PaymentProcessor, amount int) {
	err := p.Pay(amount)
	if err != nil {
		fmt.Println("An error has been occurred. The payment was declined.")
	}
}

func sixthTask() {
	fmt.Println("\nTask 3.3.6:")

	card := CardProcessor{}
	cash := CashProcessor{}
	checkout(card, 98)
	checkout(cash, 100)
}

// 3.3.7
type Logger interface {
	Log(string)
}

type ConsoleLogger struct{}

type PrefixLogger struct {
	Prefix string
}

func (ConsoleLogger) Log(message string) { fmt.Println(message) }

func (pl PrefixLogger) Log(message string) { fmt.Println(pl.Prefix, message) }

func processOrder(logger Logger, id string) {
	logger.Log("Order has been received")
	logger.Log("Order is being processed")
	logger.Log(fmt.Sprintf("Order %s has been created", id))
}

func seventhTask() {
	fmt.Println("\nTask 3.3.7:")

	consoleLogger := ConsoleLogger{}
	processOrder(consoleLogger, "oei17")

	prefLogger := PrefixLogger{"Worldwide humus centre:"}
	processOrder(prefLogger, "oei122")
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
