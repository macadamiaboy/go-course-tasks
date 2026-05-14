package service

import (
	"fmt"

	"github.com/go-course/clean-arch-task02/domain"
	"github.com/go-course/clean-arch-task02/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(name string, price float64, stock int) (domain.Product, error) {
	p := domain.Product{Name: name, Price: price, Stock: stock}
	if err := p.Validate(); err != nil {
		return domain.Product{}, fmt.Errorf("создание товара: %w", err)
	}
	return s.repo.Save(p)
}

func (s *ProductService) List() ([]domain.Product, error) {
	return s.repo.FindAll()
}

func (s *ProductService) Buy(productID int, quantity int) error {
	p, err := s.repo.FindByID(productID)
	if err != nil {
		return err
	}
	if p.Stock < quantity {
		return fmt.Errorf("недостаточно товара на складе")
	}
	return s.repo.UpdateStock(productID, p.Stock-quantity)
}
