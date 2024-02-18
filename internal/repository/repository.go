package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jfilipedias/rinha-backend-2024-q1-go/internal/entity"
)

type CreateTransactionResult struct {
	Balance int `json:"saldo"`
	Limit   int `json:"limite"`
}

var ErrCustomerNotFound = errors.New("customer not found")
var ErrInsufficientLimit = errors.New("insufficient limit")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindStatementByCustomerID(customerID int) (*entity.Statement, error) {
	sql := `
		SELECT c.debit_limit AS limit, b.value AS total, NOW() AS created_at 
		FROM customers AS c 
		LEFT JOIN balances AS b ON b.customer_id = c.id
		WHERE c.id = $1
	`

	var balance entity.Balance
	err := r.db.QueryRow(context.Background(), sql, customerID).Scan(&balance.Limit, &balance.Total, &balance.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCustomerNotFound
		}

		return nil, err
	}

	rows, err := r.db.Query(context.Background(),
		"SELECT value, type, description, created_at FROM transactions WHERE customer_id = $1 ORDER BY created_at DESC LIMIT 10",
		customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]entity.Transaction, 0)
	for rows.Next() {
		var transaction entity.Transaction
		err := rows.Scan(&transaction.Value, &transaction.Type, &transaction.Description, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	statement := &entity.Statement{
		Balance:            balance,
		LatestTransactions: transactions,
	}

	return statement, nil
}

func (r *Repository) CreateTransaction(transaction *entity.Transaction) (*CreateTransactionResult, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	sql := `
		SELECT c.debit_limit, b.value
			FROM customers AS c 
			INNER JOIN balances AS b
			ON b.customer_id = c.id
			WHERE c.id = $1
			FOR UPDATE;
	`

	var balance entity.Balance
	err = tx.QueryRow(ctx, sql, transaction.CustomerID).Scan(&balance.Limit, &balance.Total)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCustomerNotFound
		}

		return nil, err
	}

	signedValue := transaction.Value
	if transaction.Type == "d" {
		if balance.Total-transaction.Value < -balance.Limit {
			return nil, ErrInsufficientLimit
		}
		signedValue = -transaction.Value
	}

	sql = `
		UPDATE balances
			SET value = value + $1
			WHERE customer_id = $2;
	`
	_, err = tx.Exec(ctx, sql, signedValue, transaction.CustomerID)
	if err != nil {
		return nil, err
	}

	sql = `
		INSERT INTO transactions (value, type, description, customer_id)
			VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, sql, transaction.Value, transaction.Type, transaction.Description, transaction.CustomerID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	result := &CreateTransactionResult{
		Balance: balance.Total + transaction.Value,
		Limit:   balance.Limit,
	}

	return result, nil
}
