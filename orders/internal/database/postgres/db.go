package postgresdb

import (
	"database/sql"
	"fmt"
	"time"

	"orders/internal/config"
	errors "orders/internal/errors"
	myLog "orders/internal/logger"
	"orders/internal/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Connection *sql.DB
}

func NewPostgres(cfg config.Config, logger myLog.MyLogger) *Postgres {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort, cfg.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return nil
	}
	time.Sleep(time.Minute)
	err = db.Ping()
	if err != nil {
		logger.Fatalf("Failed to ping PostgreSQL: %v", err)
		return nil
	} else {
		logger.Debugf("ping success")
	}
	time.Sleep(time.Minute)

	return &Postgres{
		Connection: db,
	}
}

func (db *Postgres) CreateOrder(order models.Order) (uuid.UUID, error) {
	var existingID uuid.UUID
	checkQuery := `
        SELECT uuid FROM orders 
        WHERE customer_id = $1 AND status = $2`

	err := db.Connection.QueryRow(checkQuery, order.CustomerID, order.Status).Scan(&existingID)
	if err == nil {
		return uuid.Nil, errors.ErrOrderAlreadyExists
	} else if err != sql.ErrNoRows {
		return uuid.Nil, errors.ErrCreateOrderInternal
	}

	query := `
        INSERT INTO orders (customer_id, total_amount, status)
        VALUES ($1, $2, $3)
        RETURNING uuid`

	var id uuid.UUID
	err = db.Connection.QueryRow(query, order.CustomerID, order.TotalAmount, order.Status).Scan(&id)
	if err != nil {
		return uuid.Nil, errors.ErrCreateOrderInternal
	}
	return id, nil
}

func (db *Postgres) UpdateOrderStatusToReady(orderID uuid.UUID, customerID uuid.UUID) error {
	var currentStatus string
	var orderCustomerID uuid.UUID
	checkQuery := `
        SELECT status, customer_id FROM orders 
        WHERE uuid = $1`

	err := db.Connection.QueryRow(checkQuery, orderID).Scan(&currentStatus, &orderCustomerID)
	if err == sql.ErrNoRows {
		return errors.ErrOrderNotFound
	} else if err != nil {
		return errors.ErrUpdateOrderInternal
	}

	if orderCustomerID != customerID {
		return errors.ErrUnauthorized
	}

	if currentStatus != "Создан" {
		return errors.ErrInvalidOrderStatus
	}

	query := `
        UPDATE orders 
        SET status = 'Готов к получению' 
        WHERE uuid = $1`

	_, err = db.Connection.Exec(query, orderID)
	if err != nil {
		return errors.ErrUpdateOrderInternal
	}
	return nil
}

func (db *Postgres) UpdateOrderStatusToReceived(orderID uuid.UUID, customerID uuid.UUID) error {
	var currentStatus string
	var orderCustomerID uuid.UUID
	checkQuery := `
        SELECT status, customer_id FROM orders 
        WHERE uuid = $1`

	err := db.Connection.QueryRow(checkQuery, orderID).Scan(&currentStatus, &orderCustomerID)
	if err == sql.ErrNoRows {
		return errors.ErrOrderNotFound
	} else if err != nil {
		return errors.ErrUpdateOrderInternal
	}

	if orderCustomerID != customerID {
		return errors.ErrUnauthorized
	}

	if currentStatus != "Готов к получению" {
		return errors.ErrInvalidOrderStatus
	}

	query := `
        UPDATE orders 
        SET status = 'Получен' 
        WHERE uuid = $1`

	_, err = db.Connection.Exec(query, orderID)
	if err != nil {
		return errors.ErrUpdateOrderInternal
	}
	return nil
}

func (db *Postgres) UpdateOrderItemQuantity(orderItemID uuid.UUID, newQuantity int, customerID uuid.UUID) error {
	var orderID uuid.UUID
	var currentStatus string
	var orderCustomerID uuid.UUID

	checkQuery := `
        SELECT order_id FROM order_items 
        WHERE uuid = $1`

	err := db.Connection.QueryRow(checkQuery, orderItemID).Scan(&orderID)
	if err == sql.ErrNoRows {
		return errors.ErrOrderItemNotFound
	} else if err != nil {
		return errors.ErrUpdateOrderItemInternal
	}

	statusQuery := `
        SELECT status, customer_id FROM orders 
        WHERE uuid = $1`

	err = db.Connection.QueryRow(statusQuery, orderID).Scan(&currentStatus, &orderCustomerID)
	if err == sql.ErrNoRows {
		return errors.ErrOrderNotFound
	} else if err != nil {
		return errors.ErrUpdateOrderInternal
	}

	if orderCustomerID != customerID {
		return errors.ErrUnauthorized
	}

	if currentStatus != "Создан" {
		return errors.ErrInvalidOrderStatus
	}

	query := `
        UPDATE order_items 
        SET quantity = $1 
        WHERE uuid = $2`

	_, err = db.Connection.Exec(query, newQuantity, orderItemID)
	if err != nil {
		return errors.ErrUpdateOrderItemInternal
	}
	return nil
}

func (db *Postgres) DeleteOrder(orderID uuid.UUID, customerID uuid.UUID) error {
	var orderCustomerID uuid.UUID
	checkQuery := `
        SELECT customer_id FROM orders 
        WHERE uuid = $1`

	err := db.Connection.QueryRow(checkQuery, orderID).Scan(&orderCustomerID)
	if err == sql.ErrNoRows {
		return errors.ErrOrderNotFound
	} else if err != nil {
		return errors.ErrDeleteOrderInternal
	}

	if orderCustomerID != customerID {
		return errors.ErrUnauthorized
	}

	query := `
        DELETE FROM orders 
        WHERE uuid = $1`

	_, err = db.Connection.Exec(query, orderID)
	if err != nil {
		return errors.ErrDeleteOrderInternal
	}
	return nil
}

func (db *Postgres) GetListOrders(status string, limit, offset int, customerID uuid.UUID) ([]models.Order, error) {
	query := `
        SELECT uuid, customer_id, total_amount, status 
        FROM orders 
        WHERE status = $1 AND customer_id = $2 
        LIMIT $3 OFFSET $4`

	rows, err := db.Connection.Query(query, status, customerID, limit, offset)
	if err != nil {
		return nil, errors.ErrGetListOrdersInternal
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.UUID, &order.CustomerID, &order.TotalAmount, &order.Status); err != nil {
			return nil, errors.ErrGetListOrdersInternal
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrGetListOrdersInternal
	}

	return orders, nil
}

func (db *Postgres) GetOrderByID(orderID uuid.UUID, customerID uuid.UUID) (models.Order, error) {
	var order models.Order
	query := `
        SELECT uuid, customer_id, order_date, total_amount, status
        FROM orders
        WHERE uuid = $1 AND customer_id = $2`

	err := db.Connection.QueryRow(query, orderID, customerID).Scan(
		&order.UUID,
		&order.CustomerID,
		&order.OrderDate,
		&order.TotalAmount,
		&order.Status,
	)
	if err == sql.ErrNoRows {
		return order, errors.ErrOrderNotFound
	} else if err != nil {
		return order, errors.ErrGetOrderInternal
	}
	return order, nil
}

func (db *Postgres) GetOrderItemsByOrderID(orderID uuid.UUID, customerID uuid.UUID) ([]models.OrderItem, error) {
	var orderCustomerID uuid.UUID
	checkQuery := `
        SELECT customer_id FROM orders 
        WHERE uuid = $1`

	err := db.Connection.QueryRow(checkQuery, orderID).Scan(&orderCustomerID)
	if err == sql.ErrNoRows {
		return nil, errors.ErrOrderNotFound
	} else if err != nil {
		return nil, errors.ErrGetOrderItemsInternal
	}

	if orderCustomerID != customerID {
		return nil, errors.ErrUnauthorized
	}

	query := `
        SELECT uuid, order_id, good_uuid, quantity
        FROM order_items
        WHERE order_id = $1`

	rows, err := db.Connection.Query(query, orderID)
	if err != nil {
		return nil, errors.ErrGetOrderItemsInternal
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.UUID,
			&item.OrderID,
			&item.GoodUUID,
			&item.Quantity,
		)
		if err != nil {
			return nil, errors.ErrGetOrderItemsInternal
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.ErrGetOrderItemsInternal
	}

	return items, nil
}
