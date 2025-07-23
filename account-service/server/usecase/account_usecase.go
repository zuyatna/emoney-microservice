package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zuyatna/emoney-microservice/account-service/server/domain"
	"github.com/zuyatna/emoney-microservice/account-service/server/repository"
	"golang.org/x/crypto/bcrypt"
)

type AccountUseCase interface {
	CreateAccount(ctx context.Context, name, email, password string) (string, error)
	LoginAccount(ctx context.Context, email, password string) (string, error)
	GetAccountByID(ctx context.Context, id string) (*domain.Account, error)
}

type accountUseCase struct {
	repo       repository.AccountRepository
	publisher  domain.AccountPublisher
	jwtSecret  []byte
	jwtExpires time.Duration
}

func NewAccountUseCase(repo repository.AccountRepository, publisher domain.AccountPublisher, jwtSecret []byte, jwtExpires time.Duration) AccountUseCase {
	return &accountUseCase{
		repo:       repo,
		publisher:  publisher,
		jwtSecret:  jwtSecret,
		jwtExpires: jwtExpires,
	}
}

func (a *accountUseCase) CreateAccount(ctx context.Context, name string, email string, password string) (string, error) {
	existing, err := a.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to check existing account: %w", err)
	}
	if existing != nil {
		return "", fmt.Errorf("account with email %s already exists", email)
	}

	account := &domain.Account{
		Name:     name,
		Email:    email,
		Password: password,
	}
	if err := a.repo.CreateAccount(ctx, account); err != nil {
		return "", fmt.Errorf("failed to create account: %w", err)
	}

	// Publish event after successful account creation
	if err := a.publisher.PublishAccountCreated(ctx, account.ID, account.Name, account.Email); err != nil {
		return "", fmt.Errorf("failed to publish account created event: %w", err)
	}

	return account.ID, nil
}

func (a *accountUseCase) LoginAccount(ctx context.Context, email string, password string) (string, error) {
	account, err := a.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get account by email: %w", err)
	}
	if account == nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	claims := &domain.CustomClaim{
		ID:    account.ID,
		Email: account.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.jwtExpires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "account-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}

func (a *accountUseCase) GetAccountByID(ctx context.Context, id string) (*domain.Account, error) {
	return a.repo.GetAccountByID(ctx, id)
}
