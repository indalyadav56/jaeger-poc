package repository

import (
	"common-service/pkg/trace"
	"context"
	"database/sql"
	"log/slog"

	"github.com/brianvoe/gofakeit/v6"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "UserRepository.CreateUser")
	defer span.End()

	_, err := r.db.ExecContext(ctx, "INSERT INTO users (email, password_hash, first_name, middle_name, last_name) VALUES ($1, $2, $3, $4, $5)", gofakeit.Email(), "password", "John", "Doe", "test")
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return err
	}

	return nil
}
