package buyersservice

import (
	"buyers-service/internal/models/request"
	"buyers-service/internal/models/response"
	"buyers-service/internal/storages/postgres"

	"context"

	"github.com/google/uuid"
)

type BuyersStorage interface {
	CreateUser(ctx context.Context, email, passwordHash string, buyer request.BuyerCreate) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (postgres.User, error)
	GetBuyer(ctx context.Context, id uuid.UUID) (response.BuyerInfo, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
