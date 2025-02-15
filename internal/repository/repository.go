package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, username, hashedPassword string) error {
	_, err := r.db.Exec(ctx, QueryCreateUser, username, hashedPassword)
	return err
}

func (r *Repository) GetUserByLogin(ctx context.Context, login string) (string, string, error) {
	var userID, passwordHash string
	err := r.db.QueryRow(ctx, QueryGetUserByLogin, login).Scan(&userID, &passwordHash)
	return userID, passwordHash, err
}

func (r *Repository) GetUserCoins(ctx context.Context, userID string) (int32, error) {

	// Выполняем запрос для получения баланса пользователя
	var coins int32
	err := r.db.QueryRow(ctx, getUserCoinsQuery, userID).Scan(&coins)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("error querying user coins: %v", err)
	}

	return coins, nil
}

func (r *Repository) GetCoinHistory(ctx context.Context, userID string) ([]struct {
	FromUser string
	ToUser   string
	Amount   int32
}, error) {
	rows, err := r.db.Query(ctx, getCoinHistoryQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching coin history: %v", err)
	}
	defer rows.Close()

	transactions := []struct {
		FromUser string
		ToUser   string
		Amount   int32
	}{} // Явная инициализация пустого среза

	for rows.Next() {
		var fromUser, toUser string
		var amount int32
		if err := rows.Scan(&fromUser, &toUser, &amount); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		transactions = append(transactions, struct {
			FromUser string
			ToUser   string
			Amount   int32
		}{
			FromUser: fromUser,
			ToUser:   toUser,
			Amount:   amount,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return transactions, nil // Если транзакций не было, вернётся пустой массив, а не null
}

func (r *Repository) GetUserInventory(ctx context.Context, userID string) ([]struct {
	Quantity int
	Type     string
}, error) {
	rows, err := r.db.Query(ctx, "SELECT type, quantity FROM inventories WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching inventory: %v", err)
	}
	defer rows.Close()

	var inventory []struct {
		Quantity int
		Type     string
	}

	for rows.Next() {
		var itemType string
		var quantity int
		if err := rows.Scan(&itemType, &quantity); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		inventory = append(inventory, struct {
			Quantity int
			Type     string
		}{
			Quantity: quantity,
			Type:     itemType,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return inventory, nil // Возвращаем инвентарь
}

func (r *Repository) SubtractUserCoins(ctx context.Context, user_id string, amount int) (int, error) {
	var updatedCoins int
	err := r.db.QueryRow(ctx, QuerySubtractUserCoins, amount, user_id).Scan(&updatedCoins)
	if err != nil {
		log.Printf("[ERROR] Failed to subtract coins for user %s: %v", user_id, err)
		return 0, fmt.Errorf("error updating coins: %v", err)
	}
	return updatedCoins, nil
}

func (r *Repository) AddItemToInventory(ctx context.Context, user_id, itemType string, quantity int) error {
	// Выполнение запроса на добавление или обновление инвентаря, используя запрос из sql_queries.go
	_, err := r.db.Exec(ctx, QueryAddOrUpdateItemInInventory, user_id, itemType, quantity)

	if err != nil {
		log.Printf("[ERROR] Failed to update inventory for user %s: %v", user_id, err)
		return fmt.Errorf("error updating inventory: %v", err)
	}

	// Логируем успешное добавление
	log.Printf("[INFO] Added item to inventory for user %s: type=%s, quantity=%d", user_id, itemType, quantity)

	return nil
}

func (r *Repository) TransferCoins(ctx context.Context, fromUserID, toUserID string, amount int32) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Списываем монеты у отправителя
	var senderNewBalance int
	err = tx.QueryRow(ctx, QuerySubtractUserCoins, amount, fromUserID).Scan(&senderNewBalance)
	if err != nil {
		return fmt.Errorf("failed to deduct coins from sender: %v", err)
	}

	// Зачисляем монеты получателю
	var recipientNewBalance int
	err = tx.QueryRow(ctx, QueryAddUserCoins, amount, toUserID).Scan(&recipientNewBalance)
	if err != nil {
		return fmt.Errorf("failed to add coins to recipient: %v", err)
	}

	// Записываем транзакцию
	_, err = tx.Exec(ctx, QueryInsertTransaction, fromUserID, toUserID, amount)
	if err != nil {
		return fmt.Errorf("failed to insert transaction record: %v", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("[INFO] Transferred %d coins from %s to %s", amount, fromUserID, toUserID)
	return nil
}
