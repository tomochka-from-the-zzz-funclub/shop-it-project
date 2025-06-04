package auth

import (
	"context"
	"errors"
	postgresdb "market/internal/database/postgres"
	"market/internal/models/request"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo      UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

type UserRepository interface {
	CreateSeller(ctx context.Context, email, passwordHash string, seller request.Seller) (uuid.UUID, error)
	GetSellerByEmail(ctx context.Context, email string) (postgresdb.Seller, error)
}

func NewAuthService(repo UserRepository, jwtSecret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string, seller request.Seller) (uuid.UUID, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, err
	}
	return s.repo.CreateSeller(ctx, email, string(hashedPassword), seller)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetSellerByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	claims := Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}
	return uuid.UUID{}, errors.New("invalid token")
}
