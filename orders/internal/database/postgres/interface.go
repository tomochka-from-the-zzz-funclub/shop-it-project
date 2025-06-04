package postgresdb

import (
	"orders/internal/models"

	"github.com/google/uuid"
)

type InterfacePostgresDB interface {
	CreateOrder(order models.Order) (uuid.UUID, error)
	UpdateOrderStatusToReady(orderID uuid.UUID, customerID uuid.UUID) error
	UpdateOrderStatusToReceived(orderID uuid.UUID, customerID uuid.UUID) error
	UpdateOrderItemQuantity(orderItemID uuid.UUID, newQuantity int, customerID uuid.UUID) error
	DeleteOrder(orderID uuid.UUID, customerID uuid.UUID) error
	GetListOrders(status string, limit, offset int, customerID uuid.UUID) ([]models.Order, error)
	GetOrderByID(orderID uuid.UUID, customerID uuid.UUID) (models.Order, error)
	GetOrderItemsByOrderID(orderID uuid.UUID, customerID uuid.UUID) ([]models.OrderItem, error)
	//GetGoodCardByID(goodUUID uuid.UUID) (models.GoodCard, error)

	//bag
	AddGoodToBag(bagID uuid.UUID, goodID uuid.UUID, customerID uuid.UUID) error
	RemoveGoodFromBag(bagID uuid.UUID, goodID uuid.UUID, customerID uuid.UUID) error
	CheckoutBag(bagID uuid.UUID, customerID uuid.UUID) (models.Order, error)
	//GetGoodsFromBag(bagUUID uuid.UUID, customerID uuid.UUID) ([]models.Good, error)
	GetBagIDByUserID(customerID uuid.UUID) (uuid.UUID, error)
}
