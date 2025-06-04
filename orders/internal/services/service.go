package services

import (
	config "orders/internal/config"
	database "orders/internal/database/postgres"
	myErrors "orders/internal/errors"
	myLog "orders/internal/logger"
	"orders/internal/models"

	"github.com/google/uuid"
)

type Srv struct {
	db database.InterfacePostgresDB
}

func NewSrv(cfg config.Config, logger myLog.MyLogger) *Srv {
	base := database.NewPostgres(cfg, logger)
	return &Srv{
		db: base,
	}
}

// SrvGetListOrders получает список заказов по статусу с учетом лимита и смещения
func (srv *Srv) SrvGetListOrders(status string, limit, offset int, customerID uuid.UUID) ([]models.Order, error) {
	// Рассчитываем начальную позицию
	startPosition := offset * limit

	// Вызываем метод базы данных для получения списка заказов
	return srv.db.GetListOrders(status, limit, startPosition, customerID)
}

// SrvCreateOrder реализует создание заказа
func (s *Srv) SrvCreateOrder(order models.Order) (uuid.UUID, error) {
	return s.db.CreateOrder(order)
}

func (s *Srv) SrvGetOrderByID(orderID uuid.UUID, customerID uuid.UUID) (models.Order, []models.Good, error) {
	// Получаем заказ из базы данных
	order, err := s.db.GetOrderByID(orderID, customerID)
	if err != nil {
		// Обработка ошибок, возвращаемые из метода базы данных
		if err == myErrors.ErrGetOrderInternal {
			return models.Order{}, nil, myErrors.ErrGetOrderInternal
		}
		return models.Order{}, nil, myErrors.ErrOrderNotFound // В случае, если заказ не найден
	}

	// Получаем позиции заказа
	orderItems, err := s.db.GetOrderItemsByOrderID(orderID, customerID)
	if err != nil {
		return models.Order{}, nil, myErrors.ErrGetOrderInternal // Обработка ошибок при получении позиций
	}

	// Создаем срез товаров на основе позиций заказа
	goods := make([]models.Good, len(orderItems))
	for i, item := range orderItems {
		// goodCard, err := s.db.GetGoodCardByID(item.GoodUUID)
		// if err != nil {
		// 	return models.Order{}, nil, myErrors.ErrGetGoodInternal // Обработка ошибок при получении карточки товара
		// }

		goods[i] = models.Good{
			UUID: item.GoodUUID,
		}
	}

	return order, goods, nil
}

// UpdateOrderStatusToReady - метод для смены статуса заказа на "Готов к получению"
func (s *Srv) SrvUpdateOrderStatusToReady(orderID uuid.UUID, customerID uuid.UUID) error {
	err := s.db.UpdateOrderStatusToReady(orderID, customerID)
	if err != nil {
		return err // Возвращаем ошибку, если что-то пошло не так
	}
	return nil
}

// UpdateOrderStatusToReceived - метод для смены статуса заказа на "Получен"
func (s *Srv) SrvUpdateOrderStatusToReceived(orderID uuid.UUID, customerID uuid.UUID) error {
	err := s.db.UpdateOrderStatusToReceived(orderID, customerID)
	if err != nil {
		return err // Возвращаем ошибку, если что-то пошло не так
	}
	return nil
}

// UpdateOrderItemQuantity - метод для изменения количества товара в позиции заказа
func (s *Srv) SrvUpdateOrderItemQuantity(orderItemID uuid.UUID, newQuantity int, customerID uuid.UUID) error {
	err := s.db.UpdateOrderItemQuantity(orderItemID, newQuantity, customerID)
	if err != nil {
		return err // Возвращаем ошибку, если что-то пошло не так
	}
	return nil
}

// DeleteOrder - метод для удаления заказа
func (s *Srv) SrvDeleteOrder(orderID uuid.UUID, customerID uuid.UUID) error {
	err := s.db.DeleteOrder(orderID, customerID)
	if err != nil {
		return err // Возвращаем ошибку, если что-то пошло не так
	}
	return nil
}
