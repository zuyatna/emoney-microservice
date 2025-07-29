package repository

import (
	"context"
	"github.com/zuyatna/emoney-microservice/transaction-service/server/model"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *model.Transaction) error
	FindHistoryByAccountID(ctx context.Context, accountID string, page, limit int) ([]*model.Transaction, int64, error)
	CreateAccount(ctx context.Context, acc *model.Account) error
}
