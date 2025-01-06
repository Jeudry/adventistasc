package store

import (
	"context"
	"database/sql"
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
