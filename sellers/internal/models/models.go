package models

// Seller представляет таблицу sellers
type Seller struct {
	Name             string `json:"name"`
	Passport         int    `json:"passport"`
	Telephone_Number string `json:"telephone_number"`
	Description      string `json:"description"`
}
