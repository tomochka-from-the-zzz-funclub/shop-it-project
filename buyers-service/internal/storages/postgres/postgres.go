package postgres

import (
	"buyers-service/internal/config"
	"buyers-service/internal/errors"
	"buyers-service/internal/logger"
	"buyers-service/internal/models/request"
	"buyers-service/internal/models/response"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
	lg *logger.Logger
}

func NewPostgres(cfg config.PostgresCFG, lg *logger.Logger) PostgresStorage {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Db)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		lg.LogErrorf("Failed to open database: %v", err)
		panic(fmt.Errorf("failed to open database: %w", err))
	}
	if err := db.Ping(); err != nil {
		db.Close()
		lg.LogErrorf("Failed to ping database: %v", err)
		panic(fmt.Errorf("failed to ping database: %w", err))
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	lg.LogInfof("Successfully connected to PostgreSQL at %s:%d/%s", cfg.Host, cfg.Port, cfg.Db)
	return PostgresStorage{db: db, lg: lg}
}

func (p *PostgresStorage) Close() {
	if err := p.db.Close(); err != nil {
		p.lg.LogErrorf("Failed to close database connection: %v", err)
	} else {
		p.lg.LogInfo("Database connection closed")
	}
}

func (p *PostgresStorage) CreateUser(ctx context.Context, email, passwordHash string, buyer request.BuyerCreate) (uuid.UUID, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.lg.LogErrorf("Failed to start transaction: %v", err)
		return uuid.UUID{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var userID uuid.UUID
	queryUser := `
        INSERT INTO users (email, password_hash)
        VALUES ($1, $2)
        RETURNING id`
	err = tx.QueryRowContext(ctx, queryUser, email, passwordHash).Scan(&userID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return uuid.UUID{}, errors.ErrDuplicateEmail
		}
		p.lg.LogErrorf("Failed to create user with email %s: %v", email, err)
		return uuid.UUID{}, fmt.Errorf("failed to create user: %w", err)
	}

	queryBuyer := `
        INSERT INTO buyer_info (id, name, phone_number, gender, birthdate)
        VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, queryBuyer, userID, buyer.Name, buyer.PhoneNumber, buyer.Gender, buyer.Birthdate)
	if err != nil {
		p.lg.LogErrorf("Failed to create buyer for user %s: %v", userID, err)
		return uuid.UUID{}, fmt.Errorf("failed to create buyer: %w", err)
	}

	if err = tx.Commit(); err != nil {
		p.lg.LogErrorf("Failed to commit transaction: %v", err)
		return uuid.UUID{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	p.lg.LogInfof("Successfully created user %s and buyer with same id", userID)
	return userID, nil
}

func (p *PostgresStorage) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	row := p.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			p.lg.LogErrorf("User with email %s not found", email)
			return User{}, errors.ErrBuyerNotFound
		}
		p.lg.LogErrorf("Failed to get user with email %s: %v", email, err)
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (p *PostgresStorage) GetBuyer(ctx context.Context, id uuid.UUID) (response.BuyerInfo, error) {
	var buyer response.BuyerInfo
	query := `SELECT id, name, phone_number, gender, birthdate FROM buyer_info WHERE id = $1`
	row := p.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&buyer.ID, &buyer.Name, &buyer.PhoneNumber, &buyer.Gender, &buyer.Birthdate)
	if err != nil {
		if err == sql.ErrNoRows {
			p.lg.LogErrorf("Buyer with ID %s not found", id)
			return response.BuyerInfo{}, errors.ErrBuyerNotFound
		}
		p.lg.LogErrorf("Failed to get buyer with ID %s: %v", id, err)
		return response.BuyerInfo{}, fmt.Errorf("failed to get buyer: %w", err)
	}
	return buyer, nil
}

func (p *PostgresStorage) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		p.lg.LogErrorf("Failed to delete user with ID %s: %v", id, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		p.lg.LogErrorf("Failed to check rows affected for user %s: %v", id, err)
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rows == 0 {
		p.lg.LogErrorf("User with ID %s not found", id)
		return errors.ErrBuyerNotFound
	}
	p.lg.LogInfof("Successfully deleted user %s", id)
	return nil
}
