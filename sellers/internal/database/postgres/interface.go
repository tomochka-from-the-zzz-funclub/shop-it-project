package postgresdb

import (
	"market/internal/models"

	"github.com/google/uuid"
)

type InterfacePostgresDB interface {
	CreateSeller(seller models.Seller) (uuid.UUID, error)

	DeleteSeller(sellerID uuid.UUID) error

	UpdateSeller(uuid uuid.UUID, seller models.Seller) error
}
