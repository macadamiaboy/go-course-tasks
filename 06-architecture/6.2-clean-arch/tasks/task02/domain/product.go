package domain

import "errors"

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
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
