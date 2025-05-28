package services

import (
	"market/internal/models"

	"github.com/google/uuid"
)

type InterfaceService interface {
	Create(s models.Seller) (uuid.UUID, error)
	Update(id uuid.UUID, new models.Seller) error
	Delete(id uuid.UUID) error
}
