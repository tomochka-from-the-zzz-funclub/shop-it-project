package buyersservice

import (
	"buyers-service/internal/logger"
	"buyers-service/internal/models/request"
	"buyers-service/internal/models/response"
	"buyers-service/internal/services/auth"
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUnauthorized = errors.New("unauthorized access")
	ErrBuyerExists  = errors.New("buyer already exists for this user")
)

type BuyersService struct {
	db          BuyersStorage
	authService *auth.AuthService
	lg          *logger.Logger
}

func NewBuyersService(db BuyersStorage, authService *auth.AuthService, lg *logger.Logger) BuyersService {
	return BuyersService{
		db:          db,
		authService: authService,
		lg:          lg,
	}
}

func (b *BuyersService) Register(ctx context.Context, email, password string, buyer request.BuyerCreate) (uuid.UUID, error) {
	userID, err := b.authService.Register(ctx, email, password, buyer)
	if err != nil {
		b.lg.LogErrorf("Failed to register user with email %s: %v", email, err)
		return uuid.UUID{}, err
	}
	b.lg.LogInfof("Successfully registered user %s", userID)
	return userID, nil
}

func (b *BuyersService) Login(ctx context.Context, email, password string) (string, error) {
	token, err := b.authService.Login(ctx, email, password)
	if err != nil {
		b.lg.LogErrorf("Failed to login user with email %s: %v", email, err)
		return "", err
	}
	b.lg.LogInfof("Successfully logged in user with email %s", email)
	return token, nil
}

func (b *BuyersService) GetBuyer(ctx context.Context, id uuid.UUID) (response.BuyerInfo, error) {
	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		b.lg.LogErrorf("Failed to get userID from context")
		return response.BuyerInfo{}, ErrUnauthorized
	}

	if userID != id {
		b.lg.LogErrorf("User %s attempted to access buyer %s", userID, id)
		return response.BuyerInfo{}, ErrUnauthorized
	}

	buyer, err := b.db.GetBuyer(ctx, id)
	if err != nil {
		b.lg.LogErrorf("Failed to get buyer %s: %v", id, err)
		return response.BuyerInfo{}, err
	}

	b.lg.LogInfof("Successfully retrieved buyer %s", id)
	return buyer, nil
}

func (b *BuyersService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		b.lg.LogErrorf("Failed to get userID from context")
		return ErrUnauthorized
	}

	if userID != id {
		b.lg.LogErrorf("User %s attempted to delete user %s", userID, id)
		return ErrUnauthorized
	}

	if err := b.db.DeleteUser(ctx, id); err != nil {
		b.lg.LogErrorf("Failed to delete user %s: %v", id, err)
		return err
	}

	b.lg.LogInfof("Successfully deleted user %s", id)
	return nil
}
