package services

import (
	"goods/internal/models"

	"github.com/google/uuid"
)

type InterfaceService interface {
	SrvCreateGoodCard(card models.GoodCard) (uuid.UUID, error)

	SrvDeleteGoodCard(cardId uuid.UUID) error

	SrvUpdateGoodCard(uuid uuid.UUID, card models.GoodCard) error

	SrvReadGood(uuid uuid.UUID) (models.Good, error)

	SrvAddCountGood(uuid uuid.UUID, number int) (int, error) //суммарное количество товара

	SrvDeleteCountGood(uuid uuid.UUID, number int) (int, error)

	SrvCreateGood(cardID uuid.UUID, quantity int) error
	SrvDeleteGood(goodID uuid.UUID) error
}
