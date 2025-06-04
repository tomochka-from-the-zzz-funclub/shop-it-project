package services

import (
	"orders/internal/models"

	"github.com/google/uuid"
)

// AddGoodToBag - метод для добавления товара в корзину
func (s *Srv) AddGoodToBag(bagID uuid.UUID, goodID uuid.UUID, customerID uuid.UUID) error {
	err := s.db.AddGoodToBag(bagID, goodID, customerID)
	if err != nil {
		return err // Возвращаем ошибку, если что-то пошло не так
	}
	return nil
}

// RemoveGoodFromBag - метод для удаления товара из корзины
func (s *Srv) RemoveGoodFromBag(bagID uuid.UUID, goodID uuid.UUID, customerID uuid.UUID) error {
	err := s.db.RemoveGoodFromBag(bagID, goodID, customerID)
	if err != nil {
		return err // Возвращаем ошибку, если что-то пошло не так
	}
	return nil
}

// CheckoutBag - метод для выкупа корзины
func (s *Srv) CheckoutBag(bagID uuid.UUID, customerID uuid.UUID) (models.Order, error) {
	order, err := s.db.CheckoutBag(bagID, customerID)
	if err != nil {
		return models.Order{}, err // Возвращаем ошибку, если что-то пошло не так
	}
	return order, nil
}

// // GetGoodsFromBag - метод для получения товаров из корзины
// func (s *Srv) GetGoodsFromBag(bagID uuid.UUID, customerID uuid.UUID) ([]models.Good, error) {
// 	goods, err := s.db.GetGoodsFromBag(bagID, customerID)
// 	if err != nil {
// 		return nil, err // Возвращаем ошибку, если что-то пошло не так
// 	}
// 	return goods, nil
// }

func (s *Srv) GetBagIDByUser(userId uuid.UUID) (uuid.UUID, error) {

	bagID, err := s.db.GetBagIDByUserID(userId)
	if err != nil {
		return uuid.UUID{}, err
	}
	return bagID, nil
}
