package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Account struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomClaim struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type AccountPublisher interface {
	PublishAccountCreated(ctx context.Context, id, name, email string) error
}
