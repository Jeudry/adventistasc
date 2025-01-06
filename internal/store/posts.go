package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store/models"
	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *models.PostsModel) error {
	query := `INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(
		ctx, query, post.Content, post.Title, post.UserID, pq.Array(&post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return err
}

func (s *PostsStore) RetrieveById(ctx context.Context, id int64) (*models.PostsModel, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at, version FROM posts WHERE id = $1`

	var post models.PostsModel

	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Content, &post.Title, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt, &post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostsStore) Update(ctx context.Context, postUpdated *models.PostsModel) error {
	query := `UPDATE posts SET content = $1, title = $2 WHERE id = $5 AND version = version + 1 RETURNING version`

	err := s.db.QueryRowContext(ctx, query, postUpdated.Content, postUpdated.Title, postUpdated.Version).Scan(
		&postUpdated.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return err
}

func (s *PostsStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
