package request

import "time"

type Seller struct {
	Name        string    `json:"name"`
	PhoneNumber int64     `json:"phone"`
	Gender      bool      `json:"gender"`
	Birthdate   time.Time `json:"birthdate"`
}
