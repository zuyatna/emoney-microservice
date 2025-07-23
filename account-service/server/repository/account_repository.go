package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zuyatna/emoney-microservice/account-service/server/domain"
	"golang.org/x/crypto/bcrypt"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *domain.Account) error
	GetAccountByID(ctx context.Context, id string) (*domain.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error)
}

type accountRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewAccountRepository(db *sql.DB, redis *redis.Client) AccountRepository {
	return &accountRepository{
		db:    db,
		redis: redis,
	}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	newUUID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}
	account.ID = newUUID.String()

	account.Password = string(hashedPassword)
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	query := `INSERT INTO accounts (id, name, email, password, balance, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = r.db.ExecContext(ctx, query, account.ID, account.Name, account.Email, account.Password, account.Balance, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (r *accountRepository) GetAccountByID(ctx context.Context, id string) (*domain.Account, error) {
	// Check Redis cache first
	cacheKey := fmt.Sprintf("account:%s", id)
	cachedAccount, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var account domain.Account
		if json.Unmarshal([]byte(cachedAccount), &account) == nil {
			return &account, nil
		}
	}

	query := `SELECT id, name, email, password, balance, created_at, updated_at FROM accounts WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	account := &domain.Account{}
	err = row.Scan(&account.ID, &account.Name, &account.Email, &account.Password, &account.Balance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get account by ID: %w", err)
	}

	// Store in Redis cache
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal account to JSON: %w", err)
	}
	if err := r.redis.Set(ctx, cacheKey, accountJSON, 24*time.Hour).Err(); err != nil {
		return nil, fmt.Errorf("failed to set account in Redis cache: %w", err)
	}

	return account, nil
}

func (r *accountRepository) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	query := `SELECT id, name, email, password, balance, created_at, updated_at FROM accounts WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)

	account := &domain.Account{}
	err := row.Scan(&account.ID, &account.Name, &account.Email, &account.Password, &account.Balance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get account by email: %w", err)
	}

	return account, nil
}
