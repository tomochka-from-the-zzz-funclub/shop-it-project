package transport

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *HandlersBuilder) Parse(tokenString string) (buyerID uuid.UUID, err error) {
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
