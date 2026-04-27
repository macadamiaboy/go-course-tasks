package main

import (
	"fmt"
	"time"
)

// 3.1.1
type Profile struct {
	Name     string
	Age      int
	IsActive bool
}

func firstTask() {
	fmt.Println("Task 3.1.1:")
	prof := Profile{Name: "Theo", Age: 32, IsActive: true}
	fmt.Printf("Name: %s, Age: %d, IsActive: %v\n", prof.Name, prof.Age, prof.IsActive)
}

// 3.1.2
type AppConfig struct {
	Host  string
	Port  int
	Debug bool
}

func secondTask() {
	fmt.Println("\nTask 3.1.2:")
	var cfg AppConfig
	fmt.Println(cfg)
	fmt.Printf("Host: %s, Port: %d, Debug: %v\n", cfg.Host, cfg.Port, cfg.Debug)
}

// 3.1.3
type Address struct {
	City   string
	Street string
}

type Employee struct {
	Name    string
	Address Address
}

func thirdTask() {
	fmt.Println("\nTask 3.1.3:")
	e := Employee{Name: "Mira", Address: Address{City: "Kazan", Street: "Baumana"}}
	fmt.Printf("%s: %s, %s\n", e.Name, e.Address.City, e.Address.Street)
}

// 3.1.4
type Package struct {
	ID     string
	Weight int
}

type Destination struct {
	City string
	Zip  string
}

type Shipment struct {
	Package     Package
	Destination Destination
}

func fourthTask() {
	fmt.Println("\nTask 3.1.4:")
	shipment := Shipment{Package: Package{ID: "EI1938", Weight: 18}, Destination: Destination{City: "Makhachkala", Zip: "367016"}}
	fmt.Printf("ID: %s, Destination: %s\n", shipment.Package.ID, shipment.Destination.City)
}

// 3.1.5
type Audit struct {
	CreatedAt string
	UpdatedAt string
}

type Article struct {
	Title string
	Audit
}

func fifthTask() {
	fmt.Println("\nTask 3.1.5:")
	article := Article{Title: "Luzhniki", Audit: Audit{CreatedAt: "14 Sep 1967", UpdatedAt: "19 Oct 1968"}}
	fmt.Printf("Title: %s, Created: %s, Upd: %s\n", article.Title, article.CreatedAt, article.UpdatedAt)
}

// 3.1.6
type ContactInfo struct {
	Phone string
	Email string
}

// Address has already been created

type Client struct {
	ID      string
	Address Address
	ContactInfo
}

func sixthTask() {
	fmt.Println("\nTask 3.1.6:")
	client := Client{
		ID:          "lkedrs129",
		Address:     Address{City: "Kazan", Street: "Baumana"},
		ContactInfo: ContactInfo{Phone: "+79999999999", Email: "1@mail.ru"},
	}
	fmt.Printf("ID: %s, City: %s, Email: %s\n", client.ID, client.Address.City, client.Email)
}

// 3.1.7
type CourseEnrollment struct {
	ID      int64
	Status  string
	Student Student
	Course  Course
}

type PersonalInfo struct {
	Name    string
	Surname string
	Email   string
}

type Student struct {
	ID int64
	PersonalInfo
}

type Skill string

type Admin struct {
	ID            int64
	Qualification []Skill
	PersonalInfo
}

type Course struct {
	ID        int64
	Name      string
	StartDate time.Time
	Admin     Admin
}

func seventhTask() {
	fmt.Println("\nTask 3.1.7:")
	student := Student{ID: 123, PersonalInfo: PersonalInfo{Name: "Alexey", Surname: "Griboedov", Email: "mushroom@mail.ru"}}
	admin := Admin{
		ID:            33,
		Qualification: []Skill{"Golang", "Python", "Web3", "Solidity"},
		PersonalInfo:  PersonalInfo{Name: "Zakhar", Surname: "Ryabinin", Email: "zryab@mail.ru"},
	}
	course := Course{ID: 5, Name: "Go Junior", StartDate: time.Date(2026, time.October, 13, 12, 0, 0, 0, time.UTC), Admin: admin}
	enrollment := CourseEnrollment{ID: 189, Status: "InProcess", Student: student, Course: course}

	fmt.Printf("ID: %d, Status: %s, Student id: %d, Course: %s\n", enrollment.ID, enrollment.Status, enrollment.Student.ID, enrollment.Course.Name)
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
