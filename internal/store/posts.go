package store

import (
	"context"
	"database/sql"
)

type PostsStore struct {
	db *sql.DB
}




func (s *PostsStore) Create(ctx context.Context, post *PostModel) error {
	// Implement the Create method here
	return nil
}