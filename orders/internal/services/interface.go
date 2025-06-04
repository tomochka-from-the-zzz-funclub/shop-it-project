package services

import (
	"orders/internal/models"

	"github.com/google/uuid"
)

type InterfaceService interface { // Получить список заказов по статусу с лимитом и смещением
	// Получить список заказов по статусу с лимитом и смещением
	SrvGetListOrders(status string, limit, offset int, customerID uuid.UUID) ([]models.Order, error)

	// // Создать заказ и вернуть его UUID
	SrvCreateOrder(order models.Order) (uuid.UUID, error)

	// Получить заказ и список товаров по UUID заказа
	SrvGetOrderByID(orderID uuid.UUID, customerID uuid.UUID) (models.Order, []models.Good, error)
// admins
	// Обновить статус заказа на "Готов к получению"
	SrvUpdateOrderStatusToReady(orderID uuid.UUID, customerID uuid.UUID) error

	// Обновить статус заказа на "Получен"
	SrvUpdateOrderStatusToReceived(orderID uuid.UUID, customerID uuid.UUID) error

	// // Обновить количество товара в позиции заказа
	SrvUpdateOrderItemQuantity(orderItemID uuid.UUID, newQuantity int, customerID uuid.UUID) error

	// Удалить заказ по UUID
	SrvDeleteOrder(orderID uuid.UUID, customerID uuid.UUID) error

	AddGoodToBag(bagID uuid.UUID, goodID uuid.UUID, customerID uuid.UUID) error
	// RemoveGoodFromBag - метод для удаления товара из корзины
	RemoveGoodFromBag(bagID uuid.UUID, goodID uuid.UUID, customerID uuid.UUID) error

	// CheckoutBag - метод для выкупа корзины
	CheckoutBag(bagID uuid.UUID, customerID uuid.UUID) (models.Order, error)
	// GetGoodsFromBag - метод для получения товаров из корзины
	// GetGoodsFromBag(bagID uuid.UUID, customerID uuid.UUID) ([]models.Good, error)

	GetBagIDByUser(customerID uuid.UUID) (uuid.UUID, error)
}
