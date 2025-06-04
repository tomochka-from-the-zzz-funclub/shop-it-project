package sellersservice

import (
	"context"
	"errors"
	postgresdb "market/internal/database/postgres"
	"market/internal/logger"
	"market/internal/models/request"
	"market/internal/models/response"
	"market/internal/services/auth"

	"github.com/google/uuid"
)

var (
	ErrUnauthorized = errors.New("unauthorized access")
	ErrSellerExists = errors.New("seller already exists for this email")
)

type SellersService struct {
	db          postgresdb.InterfacePostgresDB
	authService *auth.AuthService
	lg          *logger.Logger
}

func NewSellersService(db postgresdb.InterfacePostgresDB, authService *auth.AuthService, lg *logger.Logger) *SellersService {
	return &SellersService{
		db:          db,
		authService: authService,
		lg:          lg,
	}
}

func (s *SellersService) Register(ctx context.Context, email, password string, seller request.Seller) (uuid.UUID, error) {
	sellerID, err := s.authService.Register(ctx, email, password, seller)
	if err != nil {
		s.lg.LogErrorf("Failed to register seller with email %s: %v", email, err)
		return uuid.UUID{}, err
	}
	s.lg.LogInfof("Successfully registered seller %s", sellerID)
	return sellerID, nil
}

func (s *SellersService) Login(ctx context.Context, email, password string) (string, error) {
	token, err := s.authService.Login(ctx, email, password)
	if err != nil {
		s.lg.LogErrorf("Failed to login seller with email %s: %v", email, err)
		return "", err
	}
	s.lg.LogInfof("Successfully logged in seller with email %s", email)
	return token, nil
}

func (s *SellersService) GetSeller(ctx context.Context, id uuid.UUID) (response.SellerInfo, error) {
	sellerID, ok := ctx.Value("sellerID").(uuid.UUID)
	if !ok {
		s.lg.LogErrorf("Failed to get sellerID from context")
		return response.SellerInfo{}, ErrUnauthorized
	}

	if sellerID != id {
		s.lg.LogErrorf("Seller %s attempted to access seller %s", sellerID, id)
		return response.SellerInfo{}, ErrUnauthorized
	}

	seller, err := s.db.GetSeller(ctx, id)
	if err != nil {
		s.lg.LogErrorf("Failed to get seller %s: %v", id, err)
		return response.SellerInfo{}, err
	}

	s.lg.LogInfof("Successfully retrieved seller %s", id)
	return seller, nil
}

func (s *SellersService) UpdateSeller(ctx context.Context, id uuid.UUID, seller request.Seller) error {
	sellerID, ok := ctx.Value("sellerID").(uuid.UUID)
	if !ok {
		s.lg.LogErrorf("Failed to get sellerID from context")
		return ErrUnauthorized
	}

	if sellerID != id {
		s.lg.LogErrorf("Seller %s attempted to update seller %s", sellerID, id)
		return ErrUnauthorized
	}

	if err := s.db.UpdateSeller(ctx, id, seller); err != nil {
		s.lg.LogErrorf("Failed to update seller %s: %v", id, err)
		return err
	}

	s.lg.LogInfof("Successfully updated seller %s", id)
	return nil
}

func (s *SellersService) DeleteSeller(ctx context.Context, id uuid.UUID) error {
	sellerID, ok := ctx.Value("sellerID").(uuid.UUID)
	if !ok {
		s.lg.LogErrorf("Failed to get sellerID from context")
		return ErrUnauthorized
	}

	if sellerID != id {
		s.lg.LogErrorf("Seller %s attempted to delete seller %s", sellerID, id)
		return ErrUnauthorized
	}

	if err := s.db.DeleteSeller(ctx, id); err != nil {
		s.lg.LogErrorf("Failed to delete seller %s: %v", id, err)
		return err
	}

	s.lg.LogInfof("Successfully deleted seller %s", id)
	return nil
}
