package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store/models"
)

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *models.UsersModel) error {
	query := `INSERT INTO users (user_name, first_name, last_name, email, password) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	error := s.db.QueryRowContext(
		ctx, query, user.UserName, user.FirstName, user.LastName, user.Email, user.Password,
	).Scan(&user.ID, &user.CreatedAt)

	if error != nil {
		return error
	}

	return error
}

func (s *UsersStore) RetrieveById(ctx context.Context, id int64) (*models.UsersModel, error) {
	query := `SELECT id, user_name, first_name, last_name, email, password, created_at FROM users WHERE id = $1`

	var user models.UsersModel

	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.UserName, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
