package postgresdb

import (
	"database/sql"
	"fmt"
	config "market/internal/cfg"
	myErrors "market/internal/errors"
	myLog "market/internal/logger"
	"market/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Connection *sql.DB
}

func NewPostgres(cfg config.Config) *Postgres {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort, cfg.SslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		myLog.Log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return nil //, myErrors.ErrCreatePostgresConnection
	}
	time.Sleep(time.Minute)
	err = db.Ping()
	if err != nil {
		myLog.Log.Fatalf("Failed to ping PostgreSQL: %v", err)
		return nil //, myErrors.ErrPing
	} else {
		myLog.Log.Debugf("ping success")
	}
	time.Sleep(time.Minute)

	return &Postgres{
		Connection: db,
	}
}

func (db *Postgres) CreateSeller(seller models.Seller) (uuid.UUID, error) {
	query := `
		INSERT INTO sellers (name, passport, telephone_number, description)
		VALUES ($1, $2, $3, $4)
		RETURNING uuid`
	var id uuid.UUID
	err := db.Connection.QueryRow(query, seller.Name, seller.Passport, seller.Telephone_Number, seller.Description).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (db *Postgres) DeleteSeller(sellerID uuid.UUID) error {
	var exists bool
	err := db.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM sellers WHERE uuid = $1)", sellerID).Scan(&exists)

	if err != nil {
		return myErrors.ErrDeleteSellerInternal // Если произошла ошибка при выполнении запроса
	}
	if !exists {
		return myErrors.ErrDeleteSellerNotFound
	}
	_, err_ := db.Connection.Exec(`DELETE FROM sellers WHERE uuid = $1`, sellerID)
	return err_
}

func (db *Postgres) UpdateSeller(id uuid.UUID, seller models.Seller) error {
	// Проверяем, существует ли продавец с таким ID
	var exists bool
	err := db.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM sellers WHERE uuid = $1)", id).Scan(&exists)

	if err != nil {
		return myErrors.ErrUpdateSellerInternal // Если произошла ошибка при выполнении запроса
	}
	if !exists {
		return myErrors.ErrUpdateSellerNotFound
	}

	query := "UPDATE sellers SET"
	var params []interface{}
	var setClauses []string
	paramIndex := 1 // Индекс параметра для использования в запросе

	if seller.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf(" name = $%d", paramIndex))
		params = append(params, seller.Name)
		paramIndex++
	}
	if seller.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf(" description = $%d", paramIndex))
		params = append(params, seller.Description)
		paramIndex++
	}
	if seller.Telephone_Number != "" {
		setClauses = append(setClauses, fmt.Sprintf(" telephone_number = $%d", paramIndex))
		params = append(params, seller.Telephone_Number)
		paramIndex++
	}
	if seller.Passport != 0 {
		setClauses = append(setClauses, fmt.Sprintf(" passport = $%d", paramIndex))
		params = append(params, seller.Passport)
		paramIndex++
	}

	// Проверяем, есть ли поля для обновления
	if len(setClauses) == 0 {
		return myErrors.ErrUpdateSellerNotFields
	}

	// Соединяем части запроса
	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE uuid = $%d", paramIndex)
	params = append(params, id)

	// Выполняем запрос
	_, err_ := db.Connection.Exec(query, params...)
	return err_
}
