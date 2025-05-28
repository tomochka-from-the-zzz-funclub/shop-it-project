package buyersservice

import (
	"buyers-service/internal/logger"
	"buyers-service/internal/models/request"
	"buyers-service/internal/models/response"
	"context"

	"github.com/google/uuid"
)

type BuyersService struct {
	db BuyersStorage
	lg *logger.Logger
}

func NewBuyersService(db BuyersStorage, lg *logger.Logger) BuyersService {
	return BuyersService{
		db: db,
		lg: lg,
	}
}

func (b *BuyersService) CreateBuyer(ctx context.Context, buyer request.BuyerCreate) (uuid.UUID, error) {
	return b.db.CreateBuyer(buyer)
}

func (b *BuyersService) GetBuyer(ctx context.Context, id uuid.UUID) (response.BuyerInfo, error) {
	return b.db.GetBuyer(id)
}

func (b *BuyersService) DeleteBuyer(ctx context.Context, id uuid.UUID) error {
	return b.db.DeleteBuyer(id)
}
