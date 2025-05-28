package postgresdb

import (
	"goods/internal/models"

	"github.com/google/uuid"
)

type InterfacePostgresDB interface {
	CreateGoodCard(goodCard models.GoodCard) (uuid.UUID, error)
	DeleteGoodCard(cardID uuid.UUID) error
	UpdateGoodCard(id uuid.UUID, goodCard models.GoodCard) error
	CreateGood(cardID uuid.UUID, quantity int) error
	DeleteGood(goodID uuid.UUID) error
	AddCountGood(goodID uuid.UUID, number int) (int, error)
	DeleteCountGood(goodID uuid.UUID, number int) (int, error)
	ReadGood(goodID uuid.UUID) (models.Good, error)
}
