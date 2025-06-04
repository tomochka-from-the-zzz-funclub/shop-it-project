package postgresdb

import (
	"database/sql"
	"fmt"

	errors "orders/internal/errors"
	"orders/internal/models"
	"time"

	"github.com/google/uuid"
)

func (db *Postgres) AddGoodToBag(bagID uuid.UUID, goodID uuid.UUID, customerID uuid.UUID) error {
	// Verify bag belongs to customer
	var bagCustomerID uuid.UUID
	err := db.Connection.QueryRow("SELECT customer_id FROM bag WHERE uuid = $1", bagID).Scan(&bagCustomerID)
	if err == sql.ErrNoRows {
		return errors.ErrBagNotFound
	} else if err != nil {
		return fmt.Errorf("failed to verify bag ownership: %w", err)
	}
	if bagCustomerID != customerID {
		return errors.ErrUnauthorized
	}

	// Check if good exists
	var exists bool
	err = db.Connection.QueryRow("SELECT EXISTS (SELECT 1 FROM goods WHERE uuid = $1)", goodID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check good existence: %w", err)
	}
	if !exists {
		return errors.ErrGoodNotFound
	}

	// Append goodID to Goods array
	query := `
		UPDATE bag
		SET goods = array_append(goods, $1)
		WHERE uuid = $2`
	_, err = db.Connection.Exec(query, goodID, bagID)
	if err != nil {
		return fmt.Errorf("failed to add good to bag: %w", err)
	}
	return nil
}

func (db *Postgres) RemoveGoodFromBag(bagID uuid.UUID, goodID uuid.UUID, customerID uuid.UUID) error {
	// Verify bag belongs to customer
	var bagCustomerID uuid.UUID
	err := db.Connection.QueryRow("SELECT customer_id FROM bag WHERE uuid = $1", bagID).Scan(&bagCustomerID)
	if err == sql.ErrNoRows {
		return errors.ErrBagNotFound
	} else if err != nil {
		return fmt.Errorf("failed to verify bag ownership: %w", err)
	}
	if bagCustomerID != customerID {
		return errors.ErrUnauthorized
	}

	// Remove one instance of goodID from Goods array
	query := `
		UPDATE bag
		SET goods = array_remove(goods, $1)
		WHERE uuid = $2
		RETURNING array_length(goods, 1)`
	var newLength sql.NullInt64
	err = db.Connection.QueryRow(query, goodID, bagID).Scan(&newLength)
	if err == sql.ErrNoRows {
		return errors.ErrBagNotFound
	} else if err != nil {
		return fmt.Errorf("failed to remove good from bag: %w", err)
	}
	if newLength.Valid && newLength.Int64 == 0 {
		// Check if goodID was in the array
		var exists bool
		err = db.Connection.QueryRow("SELECT EXISTS (SELECT 1 FROM unnest((SELECT goods FROM bag WHERE uuid = $1)) AS g WHERE g = $2)", bagID, goodID).Scan(&exists)
		if err == nil && !exists {
			return errors.ErrGoodNotFound
		}
	}
	return nil
}

func (db *Postgres) CheckoutBag(bagID uuid.UUID, customerID uuid.UUID) (models.Order, error) {
	// Start transaction
	tx, err := db.Connection.Begin()
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Verify bag belongs to customer
	var bagCustomerID uuid.UUID
	err = tx.QueryRow("SELECT customer_id FROM bag WHERE uuid = $1", bagID).Scan(&bagCustomerID)
	if err == sql.ErrNoRows {
		return models.Order{}, errors.ErrBagNotFound
	} else if err != nil {
		return models.Order{}, fmt.Errorf("failed to verify bag ownership: %w", err)
	}
	if bagCustomerID != customerID {
		return models.Order{}, errors.ErrUnauthorized
	}

	// Get goods from bag
	var goods []uuid.UUID
	err = tx.QueryRow("SELECT goods FROM bag WHERE uuid = $1", bagID).Scan(&goods)
	if err == sql.ErrNoRows {
		return models.Order{}, errors.ErrBagNotFound
	} else if err != nil {
		return models.Order{}, fmt.Errorf("failed to get bag goods: %w", err)
	}
	if len(goods) == 0 {
		return models.Order{}, errors.ErrBagEmpty
	}

	// Calculate quantities and total amount
	type bagItem struct {
		GoodID   uuid.UUID
		Quantity int
	}
	itemMap := make(map[uuid.UUID]int)
	for _, goodID := range goods {
		itemMap[goodID]++
	}
	var bagItems []bagItem
	var totalAmount float64
	for goodID, quantity := range itemMap {
		// Get price from goods
		var price float64
		err = tx.QueryRow("SELECT price FROM goods WHERE uuid = $1", goodID).Scan(&price)
		if err == sql.ErrNoRows {
			return models.Order{}, errors.ErrGoodNotFound
		} else if err != nil {
			return models.Order{}, fmt.Errorf("failed to get good price: %w", err)
		}
		totalAmount += price * float64(quantity)
		bagItems = append(bagItems, bagItem{GoodID: goodID, Quantity: quantity})
	}

	// Create order
	order := models.Order{
		UUID:        uuid.New(),
		CustomerID:  customerID,
		OrderDate:   time.Now(),
		TotalAmount: totalAmount,
		Status:      "Created",
	}
	_, err = tx.Exec(
		"INSERT INTO orders (uuid, customer_id, order_date, total_amount, status) VALUES ($1, $2, $3, $4, $5)",
		order.UUID, order.CustomerID, order.OrderDate, order.TotalAmount, order.Status,
	)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	for _, item := range bagItems {
		itemUUID := uuid.New()
		_, err = tx.Exec(
			"INSERT INTO order_items (uuid, order_id, good_uuid, quantity) VALUES ($1, $2, $3, $4)",
			itemUUID, order.UUID, item.GoodID, item.Quantity,
		)
		if err != nil {
			return models.Order{}, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Clear bag
	_, err = tx.Exec("UPDATE bag SET goods = '{}' WHERE uuid = $1", bagID)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to clear bag: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return models.Order{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

func (db *Postgres) GetBagIDByUserID(customerID uuid.UUID) (uuid.UUID, error) {
	var bagID uuid.UUID
	query := `SELECT uuid FROM bag WHERE customer_id = $1`
	err := db.Connection.QueryRow(query, customerID).Scan(&bagID)
	if err == sql.ErrNoRows {
		return uuid.Nil, errors.ErrBagNotFound
	} else if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get bag ID: %w", err)
	}
	return bagID, nil
}
