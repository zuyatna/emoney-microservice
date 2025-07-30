package repository

import (
	"context"
	"database/sql"
	"github.com/olivere/elastic/v7"
	"github.com/zuyatna/emoney-microservice/transaction-service/server/model"
	"log"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *model.Transaction) error
	FindHistoryByAccountID(ctx context.Context, accountID string, page, limit int) ([]*model.Transaction, int64, error)
	CreateAccount(ctx context.Context, acc *model.Account) error
}

const esIndexName = "transactions"

type transactionRepository struct {
	db *sql.DB
	es *elastic.Client
}

func NewTransactionRepository(db *sql.DB, es *elastic.Client) TransactionRepository {
	return &transactionRepository{
		db: db,
		es: es,
	}
}

func (t transactionRepository) CreateAccount(ctx context.Context, acc *model.Account) error {
	query := `INSERT INTO accounts (id, name, email, balance) VALUES ($1, $2, $3, 0) ON CONFLICT (id) DO NOTHING`
	_, err := t.db.ExecContext(ctx, query, acc.ID, acc.Name, acc.Email)
	if err != nil {
		log.Printf("Error inserting account: %v", err)
	}
	return err
}

func (t transactionRepository) CreateTransaction(ctx context.Context, tx *model.Transaction) error {
	query := `INSERT INTO transactions (id, from_account_id, to_account_id, amount, transaction_type, notes, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := t.db.ExecContext(ctx, query, tx.ID, tx.FromAccountID, tx.ToAccountID, tx.Amount, tx.TransactionType, tx.Notes, tx.CreatedAt)
	if err != nil {
		log.Printf("Error inserting transaction: %v", err)
		return err
	}

	// Index the transaction in Elasticsearch
	_, err = t.es.Index().
		Index(esIndexName).
		Id(tx.ID).
		BodyJson(tx).
		Do(ctx)
	if err != nil {
		log.Printf("Error indexing transaction in Elasticsearch: %v", err)
		return err
	}

	return nil
}

func (t transactionRepository) FindHistoryByAccountID(ctx context.Context, accountID string, page, limit int) ([]*model.Transaction, int64, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, transaction_type, notes, created_at 
			  FROM transactions 
			  WHERE from_account_id = $1 OR to_account_id = $1 
			  ORDER BY created_at DESC 
			  LIMIT $2 OFFSET $3`
	rows, err := t.db.QueryContext(ctx, query, accountID, limit, (page-1)*limit)
	if err != nil {
		log.Printf("Error querying transaction history: %v", err)
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	var transactions []*model.Transaction
	for rows.Next() {
		tx := &model.Transaction{}
		if err := rows.Scan(&tx.ID, &tx.FromAccountID, &tx.ToAccountID, &tx.Amount, &tx.TransactionType, &tx.Notes, &tx.CreatedAt); err != nil {
			log.Printf("Error scanning transaction: %v", err)
			return nil, 0, err
		}
		transactions = append(transactions, tx)
	}

	countQuery := `SELECT COUNT(*) FROM transactions WHERE from_account_id = $1 OR to_account_id = $1`
	var total int64
	err = t.db.QueryRowContext(ctx, countQuery, accountID).Scan(&total)
	if err != nil {
		log.Printf("Error counting transactions: %v", err)
		return nil, 0, err
	}

	return transactions, total, nil
}
