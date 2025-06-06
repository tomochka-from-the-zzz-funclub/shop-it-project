package models

import "github.com/google/uuid"

type GoodCard struct {
	UUID        uuid.UUID `json:"uuid"`        // Уникальный идентификатор товара (UUID)
	Price       float64   `json:"price"`       // Цена товара
	Name        string    `json:"name"`        // Название товара
	Description string    `json:"description"` // Описание товара
	Weight      float64   `json:"weight"`      // Вес товара
	SellerID    uuid.UUID `json:"sellerId"`    // Уникальный идентификатор продавца (UUID)
	IsActive    bool      `json:"isActive"`    // Статус активации товара
}

type Good struct {
	UUID     uuid.UUID `json:"uuid"`
	Card     GoodCard  `json:"card"`
	Quantity int       `json:"quantity"` // Количество товара
}
