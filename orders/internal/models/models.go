package models

import (
	"time"

	"github.com/google/uuid"
)

// Order - структура для таблицы заказов
type Order struct {
	UUID        uuid.UUID `json:"uuid"`
	CustomerID  uuid.UUID ` json:"customer_id"`
	OrderDate   time.Time `json:"order_date"`
	TotalAmount float64   `json:"total_amount" `
	Status      string    `json:"status" `
}

// OrderItem - структура для таблицы позиций заказов
type OrderItem struct {
	UUID     uuid.UUID `json:"uuid"`
	OrderID  uuid.UUID `json:"order_id"`
	GoodUUID uuid.UUID `json:"good_uuid"`
	Quantity int       `json:"quantity"`
}

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


type Bag struct {
	IdCustomer uuid.UUID
	Goods []Good
}
