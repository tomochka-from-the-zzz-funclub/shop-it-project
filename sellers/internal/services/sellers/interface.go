package sellersservice

import (
	"market/internal/models/request"
	"market/internal/models/response"

	"context"

	"github.com/google/uuid"
)

type ISellers interface {
	GetSeller(ctx context.Context, id uuid.UUID) (response.SellerInfo, error)
	UpdateSeller(ctx context.Context, id uuid.UUID, seller request.Seller) error
	DeleteSeller(ctx context.Context, id uuid.UUID) error
}
