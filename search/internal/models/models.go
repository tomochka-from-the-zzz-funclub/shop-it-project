package models

import "github.com/google/uuid"

// type Good struct {
// 	UUID        string  `json:"uuid"`
// 	Name        string  `json:"name"`
// 	Description string  `json:"description"`
// 	Price       float64 `json:"price"`
// 	Weight      float64 `json:"weight"`
// 	IsActive    bool    `json:"is_active"`
// 	SellerID    string  `json:"seller_id"`
// }

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

type SearchRequest struct {
	Query    string  `json:"query"`
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
	SortBy   string  `json:"sort_by"` // "price_asc", "price_desc", "name"
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}

type SearchResponse struct {
	Products []*Good `json:"products"`
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}
