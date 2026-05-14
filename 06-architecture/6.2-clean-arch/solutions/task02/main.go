package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-course/clean-arch-task02/domain"
	"github.com/go-course/clean-arch-task02/repository"
	"github.com/go-course/clean-arch-task02/service"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

type createRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	product, err := h.svc.Create(req.Name, req.Price, req.Stock)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, productResponse(product))
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.List()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := make([]map[string]any, 0, len(products))
	for _, p := range products {
		resp = append(resp, productResponse(p))
	}
	writeJSON(w, http.StatusOK, resp)
}

func productResponse(p domain.Product) map[string]any {
	return map[string]any{
		"id":    p.ID,
		"name":  p.Name,
		"price": p.Price,
		"stock": p.Stock,
	}
}

func main() {
	repo := repository.NewInMemoryProductRepository()
	svc := service.NewProductService(repo)
	h := NewProductHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /products", h.Create)
	mux.HandleFunc("GET /products", h.List)

	fmt.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
