package service

import (
	"errors"

	"github.com/go-course/clean-arch-task01/domain"
	"github.com/go-course/clean-arch-task01/repository"
)

// TODO: создай структуру ProductService с полем repo типа repository.ProductRepository (интерфейс!)
// TODO: создай конструктор NewProductService(repo repository.ProductRepository) *ProductService
type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// TODO: реализуй методы:
//
//	Create(name string, price float64, stock int) (domain.Product, error)
//	  - создаёт продукт через domain.Product, валидирует, сохраняет
//
//	List() ([]domain.Product, error)
//	  - возвращает все продукты
//
//	Buy(productID int, quantity int) error
//	  - находит продукт, проверяет что quantity <= Stock
//	  - если ок - уменьшает Stock через UpdateStock
//	  - если нет - возвращает ошибку "недостаточно товара на складе"
func (ps *ProductService) Create(name string, price float64, stock int) (domain.Product, error) {
	product := domain.Product{Name: name, Price: price, Stock: stock}

	err := product.Validate()
	if err != nil {
		return domain.Product{}, err
	}

	return ps.repo.Save(product)
}

func (ps *ProductService) List() ([]domain.Product, error) {
	return ps.repo.FindAll()
}

func (ps *ProductService) Buy(productID int, quantity int) error {
	product, err := ps.repo.FindByID(productID)
	if err != nil {
		return err
	}

	if product.Stock < quantity {
		return errors.New("not enough products in stock")
	}

	newQuantity := product.Stock - quantity

	err = ps.repo.UpdateStock(productID, newQuantity)
	if err != nil {
		return err
	}

	return nil
}
