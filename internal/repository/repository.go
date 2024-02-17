package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction struct {
	ID          int       `json:"-"`
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
	CustomerID  int       `json:"-"`
}

type Balance struct {
	Total     int       `json:"total"`
	Limit     int       `json:"limite"`
	CreatedAt time.Time `json:"data_extrato"`
}

type Statement struct {
	Balance            Balance       `json:"saldo"`
	LatestTransactions []Transaction `json:"ultimas_transacoes"`
}

var ErrCustomerNotFound = errors.New("customer not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindStatementByCustomerID(customerID int) (*Statement, error) {
	sql := `
		SELECT c.debit_limit AS limit, b.value AS total, NOW() AS created_at 
		FROM customers AS c 
		LEFT JOIN balances AS b ON b.customer_id = c.id
		WHERE c.id = $1
	`

	var balance Balance
	err := r.db.QueryRow(context.Background(), sql, customerID).Scan(&balance.Limit, &balance.Total, &balance.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrCustomerNotFound
		}

		return nil, err
	}

	rows, err := r.db.Query(context.Background(),
		"SELECT value, type, description, created_at FROM transactions WHERE customer_id = $1 ORDER BY created_at DESC LIMIT 1",
		customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(&transaction.Value, &transaction.Type, &transaction.Description, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	statement := &Statement{
		Balance:            balance,
		LatestTransactions: transactions,
	}

	return statement, nil
}
