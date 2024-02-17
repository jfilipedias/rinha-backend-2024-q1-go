package entity

import "time"

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
