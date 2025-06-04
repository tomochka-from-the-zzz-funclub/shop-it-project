package postgresdb

import (
	"database/sql"
	"fmt"
	config "goods/internal/cfg"
	myErrors "goods/internal/errors"
	myLog "goods/internal/logger"
	"goods/internal/models"
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

func (db *Postgres) checkSellerExists(sellerID uuid.UUID) error {
	var exists bool
	err := db.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM seller_info WHERE id = $1)", sellerID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return err
	}
	return nil
}

func (db *Postgres) CreateGoodCard(goodCard models.GoodCard, sellerID uuid.UUID) (uuid.UUID, error) {
	if err := db.checkSellerExists(sellerID); err != nil {
		return uuid.Nil, err
	}

	// Проверяем, существует ли карточка товара с полным совпадением
	var existingID uuid.UUID
	checkQuery := `
		SELECT uuid FROM good_cards 
		WHERE name = $1 AND description = $2 AND price = $3 AND weight = $4 AND seller_id = $5 AND is_active = $6`

	err := db.Connection.QueryRow(checkQuery, goodCard.Name, goodCard.Description, goodCard.Price, goodCard.Weight, goodCard.SellerID, goodCard.IsActive).Scan(&existingID)
	if err == nil {
		// Если карточка найдена, возвращаем ошибку
		return uuid.Nil, myErrors.ErrGoodCardAlreadyExists
	} else if err != sql.ErrNoRows {
		// Если произошла другая ошибка, возвращаем внутреннюю ошибку
		return uuid.Nil, myErrors.ErrCreateGoodCardInternal
	}

	// Если карточка не найдена, добавляем новую
	query := `
		INSERT INTO good_cards (price, name, description, weight, seller_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING uuid`

	var id uuid.UUID
	err = db.Connection.QueryRow(query, goodCard.Price, goodCard.Name, goodCard.Description, goodCard.Weight, goodCard.SellerID, goodCard.IsActive).Scan(&id)
	if err != nil {
		return uuid.Nil, myErrors.ErrCreateGoodCardInternal
	}
	return id, nil
}

func (db *Postgres) DeleteGoodCard(cardID uuid.UUID, sellerID uuid.UUID) error {
	if err := db.checkSellerExists(sellerID); err != nil {
		return err
	}
	var exists bool
	err := db.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM good_cards WHERE uuid = $1)", cardID).Scan(&exists)

	if err != nil {
		return myErrors.ErrDeleteGoodCardInternal // Если произошла ошибка при выполнении запроса
	}
	if !exists {
		return myErrors.ErrGoodCardNotFound // Если карточка товара не найдена
	}

	_, err = db.Connection.Exec(`DELETE FROM good_cards WHERE uuid = $1`, cardID)
	return err // Возвращаем ошибку, если она возникла, или nil, если удаление прошло успешно
}

func (db *Postgres) UpdateGoodCard(id uuid.UUID, goodCard models.GoodCard, sellerID uuid.UUID) error {
	if err := db.checkSellerExists(sellerID); err != nil {
		return err
	}
	// Проверяем, существует ли карточка товара с таким ID
	var exists bool
	err := db.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM good_cards WHERE uuid = $1)", id).Scan(&exists)

	if err != nil {
		return myErrors.ErrUpdateGoodCardInternal // Если произошла ошибка при выполнении запроса
	}
	if !exists {
		return myErrors.ErrGoodCardNotFound
	}

	query := "UPDATE good_cards SET"
	var params []interface{}
	var setClauses []string
	paramIndex := 1 // Индекс параметра для использования в запросе

	if goodCard.Price >= 0 {
		setClauses = append(setClauses, fmt.Sprintf(" price = $%d", paramIndex))
		params = append(params, goodCard.Price)
		paramIndex++
	}
	if goodCard.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf(" name = $%d", paramIndex))
		params = append(params, goodCard.Name)
		paramIndex++
	}
	if goodCard.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf(" description = $%d", paramIndex))
		params = append(params, goodCard.Description)
		paramIndex++
	}
	if goodCard.Weight >= 0 {
		setClauses = append(setClauses, fmt.Sprintf(" weight = $%d", paramIndex))
		params = append(params, goodCard.Weight)
		paramIndex++
	}
	if goodCard.SellerID != uuid.Nil {
		setClauses = append(setClauses, fmt.Sprintf(" seller_id = $%d", paramIndex))
		params = append(params, goodCard.SellerID)
		paramIndex++
	}
	setClauses = append(setClauses, fmt.Sprintf(" is_active = $%d", paramIndex))
	params = append(params, goodCard.IsActive)
	paramIndex++

	// Проверяем, есть ли поля для обновления
	if len(setClauses) == 0 {
		return myErrors.ErrUpdateGoodCardNotFields
	}

	// Соединяем части запроса
	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE uuid = $%d", paramIndex)
	params = append(params, id)

	// Выполняем запрос
	_, err_ := db.Connection.Exec(query, params...)
	return err_
}

// /////////////////////////////////////////////////good/////////////////////////////////////////////////////////////////////////////////////////////////
func (db *Postgres) CreateGood(cardID uuid.UUID, quantity int, sellerID uuid.UUID) error {
	if err := db.checkSellerExists(sellerID); err != nil {
		return err
	}
	var exists bool
	err := db.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM good_cards WHERE uuid = $1)", cardID).Scan(&exists)
	if err != nil {
		return myErrors.ErrCreateGoodInternal
	}
	if !exists {
		return myErrors.ErrGoodCardNotFound
	}

	// 1. Создаем запись в goods
	query := "INSERT INTO goods (card_id, quantity) VALUES ($1, $2)"
	_, err = db.Connection.Exec(query, cardID, quantity)
	if err != nil {
		return err
	}

	// 2. Загружаем карточку и количество
	var gc models.GoodCard
	var qty int

	err = db.Connection.QueryRow(`
		SELECT gc.uuid, gc.price, gc.name, gc.description, gc.weight, gc.seller_id, gc.is_active, g.quantity
		FROM good_cards gc
		JOIN goods g ON gc.uuid = g.card_id
		WHERE gc.uuid = $1
	`, cardID).Scan(&gc.UUID, &gc.Price, &gc.Name, &gc.Description, &gc.Weight, &gc.SellerID, &gc.IsActive, &qty)
	if err != nil {
		return err
	}

	return nil
}

func (db *Postgres) DeleteGood(goodID uuid.UUID, sellerID uuid.UUID) error {
	if err := db.checkSellerExists(sellerID); err != nil {
		return err
	}
	// Проверяем, существует ли товар с таким ID
	var exists bool
	err := db.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM goods WHERE card_id = $1)", goodID).Scan(&exists)

	if err != nil {
		return myErrors.ErrDeleteGoodInternal // Ошибка при выполнении запроса
	}
	if !exists {
		return myErrors.ErrGoodNotFound
	}

	// Удаляем товар (это также удалит связанную карточку товара из good_cards)
	query := "DELETE FROM goods WHERE card_id = $1"
	_, err = db.Connection.Exec(query, goodID)
	return err
}

func (db *Postgres) AddCountGood(goodID uuid.UUID, number int, sellerID uuid.UUID) (int, error) {
	if err := db.checkSellerExists(sellerID); err != nil {
		return 0, err
	}
	var exists bool
	err := db.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM goods WHERE card_id = $1)", goodID).Scan(&exists)

	if err != nil {
		return 0, myErrors.ErrAddCountGoodInternal
	}
	if !exists {
		return 0, myErrors.ErrGoodNotFound
	}

	query := "UPDATE goods SET quantity = quantity + $1 WHERE card_id = $2 RETURNING quantity"
	var newQuantity int
	err = db.Connection.QueryRow(query, number, goodID).Scan(&newQuantity)
	if err != nil {
		return 0, err
	}

	return newQuantity, nil
}

func (db *Postgres) DeleteCountGood(goodID uuid.UUID, number int, sellerID uuid.UUID) (int, error) {
	if err := db.checkSellerExists(sellerID); err != nil {
		return 0, err
	}
	// Проверяем, существует ли товар с таким ID
	var exists bool
	err := db.Connection.QueryRow("SELECT EXISTS(SELECT 1 FROM goods WHERE card_id = $1)", goodID).Scan(&exists)

	if err != nil {
		return 0, myErrors.ErrDeleteCountGoodInternal // Ошибка при выполнении запроса
	}
	if !exists {
		return 0, myErrors.ErrGoodNotFound
	}

	// Получаем текущее количество товара
	var currentQuantity int
	err = db.Connection.QueryRow("SELECT quantity FROM goods WHERE card_id = $1", goodID).Scan(&currentQuantity)
	if err != nil {
		return 0, err // Возвращаем ошибку, если не удалось получить текущее количество
	}

	// Проверяем, достаточно ли товара для удаления
	if currentQuantity < number {
		return currentQuantity, myErrors.ErrNotEnoughQuantity // Недостаточно товара для удаления
	}

	// Обновляем количество товара
	query := "UPDATE goods SET quantity = quantity - $1 WHERE card_id = $2 RETURNING quantity"
	var newQuantity int
	err = db.Connection.QueryRow(query, number, goodID).Scan(&newQuantity)
	if err != nil {
		return 0, err // Возвращаем ошибку, если что-то пошло не так
	}

	return newQuantity, nil // Возвращаем новое количество товара
}

func (db *Postgres) ReadGood(goodID uuid.UUID) (models.Good, error) {

	// Проверяем, существует ли товар с таким ID
	var good models.Good
	query := `
		SELECT g.uuid, g.quantity, gc.uuid, gc.name, gc.description 
		FROM goods g
		JOIN good_cards gc ON g.card_id = gc.uuid 
		WHERE g.uuid = \$1`

	err := db.Connection.QueryRow(query, goodID).Scan(&good.UUID, &good.Quantity, &good.Card.UUID, &good.Card.Name, &good.Card.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Good{}, myErrors.ErrGoodNotFound // Если товар не найден
		}
		return models.Good{}, myErrors.ErrReadCardInternal // Ошибка при выполнении запроса
	}

	return good, nil // Возвращаем структуру Good с заполненной карточкой товара
}

// func IndexGood(ctx context.Context, good models.Good) error {
// 	body, err := json.Marshal(good)
// 	if err != nil {
// 		return err
// 	}

// 	res, err := elastic.GetClient().Index(
// 		"goods",               // индекс
// 		bytes.NewReader(body), // тело запроса
// 		elastic.GetClient().Index.WithDocumentID(good.UUID.String()), // по UUID карточки
// 		elastic.GetClient().Index.WithContext(ctx),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		return fmt.Errorf("Elasticsearch index error: %s", res.String())
// 	}

// 	return nil
// }

func (db *Postgres) SearchGoods(req models.SearchRequest) (*models.SearchResponse, error) {
	offset := (req.Page - 1) * req.PageSize
	query := `
		SELECT g.card_id, gc.name, gc.description, gc.price, gc.weight, gc.is_active, g.quantity
		FROM goods g
		JOIN good_cards gc ON g.card_id = gc.uuid
		WHERE ($1 = '' OR gc.name ILIKE '%' || $1 || '%')
		AND ($2 = 0 OR gc.price >= $2)
		AND ($3 = 0 OR gc.price <= $3)
		LIMIT $4 OFFSET $5
	`
	rows, err := db.Connection.Query(query,
		req.Query, req.MinPrice, req.MaxPrice, req.PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	defer rows.Close()

	goods := []*models.Good{}
	for rows.Next() {
		var g models.Good
		var card models.GoodCard
		if err := rows.Scan(
			&g.UUID,
			&card.Name,
			&card.Description,
			&card.Price,
			&card.Weight,
			&card.IsActive,
			&g.Quantity,
		); err != nil {
			return nil, err
		}
		card.UUID = g.UUID
		g.Card = card
		goods = append(goods, &g)
	}

	// Подсчёт общего количества
	countQuery := `
		SELECT COUNT(*)
		FROM goods g
		JOIN good_cards gc ON g.card_id = gc.uuid
		WHERE ($1 = '' OR gc.name ILIKE '%' || $1 || '%')
		AND ($2 = 0 OR gc.price >= $2)
		AND ($3 = 0 OR gc.price <= $3)
	`
	var total int64
	err = db.Connection.QueryRow(countQuery, req.Query, req.MinPrice, req.MaxPrice).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &models.SearchResponse{
		Products: goods,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
