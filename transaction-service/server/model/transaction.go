package model

import "time"

type TransactionType string

const (
	Topup    TransactionType = "topup"
	Transfer TransactionType = "transfer"
)

type Transaction struct {
	ID              string
	FromAccountID   string
	ToAccountID     string
	Amount          float64
	TransactionType TransactionType
	Notes           string
	CreatedAt       time.Time
}

type Account struct {
	ID    string
	Name  string
	Email string
}
