package store

import (
	"context"
	"database/sql"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store/models"
	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *models.PostsModel) error {
	query := `INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	error := s.db.QueryRowContext(
		ctx, query, post.Content, post.Title, post.UserID, post.Tags, pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if error != nil {
		return error
	}

	return error
}
