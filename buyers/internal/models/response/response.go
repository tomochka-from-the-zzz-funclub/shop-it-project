package response

import (
	"time"

	"github.com/google/uuid"
)

type BuyerInfo struct {
	ID          uuid.UUID `json:"ID"`
	Name        string    `json:"name"`
	PhoneNumber int64     `json:"phone"`
	Gender      bool      `json:"gender"`
	Birthdate   time.Time `json:"birthdate"`
}
