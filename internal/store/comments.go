package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store/models"
	_ "github.com/lib/pq"
)

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) Create(ctx context.Context, comment *models.CommentsModel) error {
	query := `INSERT INTO comments (content, post_id, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(
		ctx, query, comment.Content, comment.PostID, comment.UserID,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		return err
	}

	return err
}

func (s *CommentsStore) RetrievePostById(ctx context.Context, id int64) (*models.CommentsModel, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.user_id, c.created_at, c.updated_at FROM comments c
		JOIN users ON comments.user_id = users.id
		WHERE c.post_id = $1
		ORDER BY c.created_at 
	`

	var comments models.CommentsModel

	err := s.db.QueryRowContext(ctx, query, id).Scan(&comments.ID, &comments.Content, &comments.PostID, &comments.UserID, &comments.CreatedAt, &comments.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &comments, nil
}
