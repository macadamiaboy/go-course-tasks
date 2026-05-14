// Задание 2: HTTP handler поверх Repository + Service
//
// Возьми код из задания 1 (domain, repository, service) и добавь HTTP-слой.
//
// Реализуй ProductHandler с двумя маршрутами:
//   POST /products        — создать товар (JSON body: {"name":"...","price":99.9,"stock":10})
//   GET  /products        — список всех товаров
//
// Требования:
//   - Handler принимает интерфейс (не конкретный тип) через конструктор
//   - При ошибке валидации — 400 Bad Request
//   - Успешное создание — 201 Created
//   - Все ответы в JSON
//
// Ожидаемый результат:
//   $ curl -s -X POST http://localhost:8080/products \
//          -H "Content-Type: application/json" \
//          -d '{"name":"iPhone 15","price":99999,"stock":10}'
//   {"id":1,"name":"iPhone 15","price":99999,"stock":10}
//
//   $ curl -s http://localhost:8080/products
//   [{"id":1,"name":"iPhone 15","price":99999,"stock":10}]
//
//   $ curl -s -X POST http://localhost:8080/products \
//          -H "Content-Type: application/json" \
//          -d '{"name":"","price":99999,"stock":10}'
//   {"error":"создание товара: название товара не может быть пустым"}
//
// Запусти: go run .

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-course/clean-arch-task02/repository"
	"github.com/go-course/clean-arch-task02/service"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// TODO: определи тип ProductHandler со полем svc *service.ProductService
// (в финальном решении — через интерфейс, но для старта используй конкретный тип)

// TODO: реализуй конструктор NewProductHandler(svc *service.ProductService) *ProductHandler

// TODO: реализуй метод (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request)
// Логика:
//   1. Декодируй JSON-тело в структуру с полями Name, Price, Stock
//   2. Вызови h.svc.Create(name, price, stock)
//   3. При ошибке — 400 {"error":"..."}
//   4. При успехе — 201 + JSON товара

// TODO: реализуй метод (h *ProductHandler) List(w http.ResponseWriter, r *http.Request)
// Логика:
//   1. Вызови h.svc.List()
//   2. Верни 200 + JSON-массив товаров

func main() {
	repo := repository.NewInMemoryProductRepository()
	svc := service.NewProductService(repo)
	_ = svc // убери после реализации
	_ = writeJSON

	mux := http.NewServeMux()

	// TODO: зарегистрируй обработчики:
	// h := NewProductHandler(svc)
	// mux.HandleFunc("POST /products", h.Create)
	// mux.HandleFunc("GET /products", h.List)

	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
