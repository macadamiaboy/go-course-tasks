package domain

import "errors"

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

func (p Product) Validate() error {
	if p.Name == "" {
		return errors.New("название товара не может быть пустым")
	}
	if p.Price <= 0 {
		return errors.New("цена должна быть больше нуля")
	}
	if p.Stock < 0 {
		return errors.New("количество на складе не может быть отрицательным")
	}
	return nil
}
