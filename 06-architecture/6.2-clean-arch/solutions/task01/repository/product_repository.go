package repository

import (
	"fmt"
	"sync"

	"github.com/go-course/clean-arch-task01/domain"
)

type ProductRepository interface {
	Save(product domain.Product) (domain.Product, error)
	FindAll() ([]domain.Product, error)
	FindByID(id int) (domain.Product, error)
	Delete(id int) error
	UpdateStock(id int, newStock int) error
}

type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products map[int]domain.Product
	counter  int
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[int]domain.Product),
	}
}

func (r *InMemoryProductRepository) Save(product domain.Product) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	product.ID = r.counter
	r.products[product.ID] = product
	return product, nil
}

func (r *InMemoryProductRepository) FindAll() ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.Product, 0, len(r.products))
	for _, p := range r.products {
		result = append(result, p)
	}
	return result, nil
}

func (r *InMemoryProductRepository) FindByID(id int) (domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.products[id]
	if !ok {
		return domain.Product{}, fmt.Errorf("товар с ID=%d не найден", id)
	}
	return p, nil
}

func (r *InMemoryProductRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[id]; !ok {
		return fmt.Errorf("товар с ID=%d не найден", id)
	}
	delete(r.products, id)
	return nil
}

func (r *InMemoryProductRepository) UpdateStock(id int, newStock int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.products[id]
	if !ok {
		return fmt.Errorf("товар с ID=%d не найден", id)
	}
	p.Stock = newStock
	r.products[id] = p
	return nil
}
