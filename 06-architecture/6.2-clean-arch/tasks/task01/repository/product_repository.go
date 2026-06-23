package repository

import (
	"errors"
	"sync"

	"github.com/go-course/clean-arch-task01/domain"
)

// TODO: создай интерфейс ProductRepository с методами:
//
//	Save(product domain.Product) (domain.Product, error)
//	FindAll() ([]domain.Product, error)
//	FindByID(id int) (domain.Product, error)
//	Delete(id int) error
//	UpdateStock(id int, newStock int) error
type ProductRepository interface {
	Save(product domain.Product) (domain.Product, error)
	FindAll() ([]domain.Product, error)
	FindByID(id int) (domain.Product, error)
	Delete(id int) error
	UpdateStock(id int, newStock int) error
}

// TODO: создай структуру InMemoryProductRepository и реализуй все методы интерфейса
// Используй map[int]domain.Product для хранения и sync.RWMutex для защиты от гонок
type InMemoryProductRepository struct {
	mu    sync.RWMutex
	data  map[int]domain.Product
	curID int
}

func (pr *InMemoryProductRepository) Save(product domain.Product) (domain.Product, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	pr.data[pr.curID] = product
	pr.curID++
	return product, nil
}

func (pr *InMemoryProductRepository) FindAll() ([]domain.Product, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	result := make([]domain.Product, 0, len(pr.data))

	for _, product := range pr.data {
		result = append(result, product)
	}

	return result, nil
}

func (pr *InMemoryProductRepository) FindByID(id int) (domain.Product, error) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	product, exists := pr.data[id]
	if !exists {
		return domain.Product{}, errors.New("not found")
	}

	return product, nil
}

func (pr *InMemoryProductRepository) Delete(id int) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	_, exists := pr.data[id]
	if !exists {
		return errors.New("not found")
	}

	delete(pr.data, id)
	return nil
}

func (pr *InMemoryProductRepository) UpdateStock(id int, newStock int) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	product, exists := pr.data[id]
	if !exists {
		return errors.New("not found")
	}

	product.Stock = newStock
	pr.data[id] = product

	return nil
}

// TODO: создай конструктор NewInMemoryProductRepository() *InMemoryProductRepository
func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{data: make(map[int]domain.Product), curID: 0}
}
