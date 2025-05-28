package entity

type SearchParams struct {
	Query           string   // текстовый запрос
	Categories      []string // фильтрация по категориям
	Brand           []string // фильтрация по брендам
	MinPrice        float64  // минимальная цена
	MaxPrice        float64  // максимальная цена
	SortBy          string   // "relevance", "price", "popularity"
	SortOrder       string   // "asc" или "desc"
	Page            int      // номер страницы (с 1)
	PageSize        int      // размер страницы
	HighlightFields []string // поля для подсветки
}
