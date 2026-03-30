package main

import "fmt"

// 3.2.1
type Book struct {
	Title string
	Pages int
}

func (b Book) Summary() string { return fmt.Sprintf("%s (%d pages)", b.Title, b.Pages) }

func firstTask() {
	fmt.Println("Task 3.2.1:")
	book := Book{Title: "The Butterfly of the Stars", Pages: 343}
	fmt.Println(book.Summary())
}

// 3.2.2
type Wallet struct {
	Balance int
}

func (w *Wallet) Deposit(amount int) { w.Balance += amount }

func secondTask() {
	fmt.Println("\nTask 3.2.2:")
	w := &Wallet{Balance: 0}
	fmt.Println("Balance:", w.Balance)

	w.Deposit(100)
	w.Deposit(474)
	fmt.Println("New balance:", w.Balance)
}

// 3.2.3
func (w Wallet) DepositCopy(amount int) { w.Balance += amount }

func thirdTask() {
	fmt.Println("\nTask 3.2.3:")
	w := &Wallet{Balance: 0}
	fmt.Println("Balance:", w.Balance)

	w.Deposit(998)
	fmt.Println("New balance:", w.Balance)

	w.DepositCopy(200)
	fmt.Println("Brand new balance:", w.Balance)
}

// 3.2.4
func resetScore(score *int) { *score = 0 }

func fourthTask() {
	fmt.Println("\nTask 3.2.4:")
	score := 123

	fmt.Println("Score:", score)
	resetScore(&score)
	fmt.Println("New score:", score)
}

// 3.2.5
type User struct {
	Name string
}

func printUserName(u *User) {
	if u == nil {
		fmt.Println("user is nil")
		return
	}
	fmt.Println(u.Name)
}

func fifthTask() {
	fmt.Println("\nTask 3.2.5:")
	var nilUser *User
	user := &User{Name: "Leo"}

	fmt.Print("Nil: ")
	printUserName(nilUser)
	fmt.Print("Not nil: ")
	printUserName(user)
}

// 3.2.6
type Timer struct {
	Seconds int
	Running bool
}

func (t *Timer) Start() { t.Running = true }

func (t *Timer) Stop() { t.Running = false }

func (t Timer) Status() string {
	if t.Running {
		return "running"
	}
	return "stopped"
}

func sixthTask() {
	fmt.Println("\nTask 3.2.6:")
	timer := Timer{Seconds: 240, Running: false}
	fmt.Println(timer.Status())

	timer.Start()
	fmt.Println(timer.Status())

	timer.Stop()
	fmt.Println(timer.Status())
}

// 3.2.7
type BankAccount struct {
	Owner   string
	Balance int
}

func (ba *BankAccount) Deposit(amount int) { ba.Balance += amount }

func (ba *BankAccount) Withdraw(amount int) bool {
	if ba.Balance < amount {
		return false
	}
	ba.Balance -= amount
	return true
}

func (ba BankAccount) Info() string {
	return fmt.Sprintf("Owner: %s, Balance: %d", ba.Owner, ba.Balance)
}

func seventhTask() {
	fmt.Println("\nTask 3.2.7:")
	acc := BankAccount{Owner: "Jack Grealish", Balance: 10000}
	fmt.Println(acc.Info())

	fmt.Println(acc.Withdraw(10001))
	fmt.Println(acc.Info())

	acc.Deposit(10)
	fmt.Println(acc.Info())

	fmt.Println(acc.Withdraw(10001))
	fmt.Println(acc.Info())
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
