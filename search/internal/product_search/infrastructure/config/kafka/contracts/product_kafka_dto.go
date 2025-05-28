package contracts

type ProductKafkaDTO struct {
	ID          string  `json:"id"`          // уникальный идентификатор
	Name        string  `json:"name"`        // название продукта
	Description string  `json:"description"` // описание продукта
	Price       float64 `json:"price"`       // цена
	Stock       int     `json:"stock"`       // количество на складе (используется как "популярность")
	Category    string  `json:"category"`    // категория продукта
	Brand       string  `json:"brand"`       // бренд продукта
}
