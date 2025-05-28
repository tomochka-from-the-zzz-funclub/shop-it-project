package buyersservice

import (
	"buyers-service/internal/models/request"
	"buyers-service/internal/models/response"

	"github.com/google/uuid"
)

type BuyersStorage interface {
	CreateBuyer(buyer request.BuyerCreate) (uuid.UUID, error)
	GetBuyer(id uuid.UUID) (buyer response.BuyerInfo, err error)
	DeleteBuyer(id uuid.UUID) (err error)
}
