package repository

import (
	"context"
	"fmt"

	"github.com/AlexCorn999/transaction-system/internal/domain"
)

// Create adds the user to the database.
func (s *Storage) Create(ctx context.Context, user domain.User) error {
	result, err := s.DB.ExecContext(ctx, "INSERT INTO users (login, password, registered_at) values ($1, $2, $3) on conflict (login) do nothing",
		user.Login, user.Password, user.RegisteredAt)
	if err != nil {
		return fmt.Errorf("postgreSQL: create %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("postgreSQL: create %w", err)
	}

	// user already exists
	if rowsAffected == 0 {
		return ErrDuplicate
	}

	return nil
}

// GetUser returns the user from the database.
func (s *Storage) GetUser(ctx context.Context, login, password string) (domain.User, error) {
	var user domain.User
	err := s.DB.QueryRowContext(ctx, "SELECT id, login, password, registered_at FROM users WHERE login=$1 AND password=$2", login, password).
		Scan(&user.ID, &user.Login, &user.Password, &user.RegisteredAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("postgreSQL: getUser %w", err)
	}
	return user, nil
}
