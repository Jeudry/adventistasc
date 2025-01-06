package store

import (
	"context"
	"database/sql"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store/models"
	_ "github.com/lib/pq"
)

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) CreatePostComment(ctx context.Context, comment *models.CommentsModel) error {
	query := `INSERT INTO comments (content, post_id, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(
		ctx, query, comment.Content, comment.PostID, comment.UserID,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		return err
	}

	return err
}

func (s *CommentsStore) RetrieveCommentsByPostId(ctx context.Context, id int64) ([]models.CommentsModel, error) {
	query := `
		SELECT c.id, c.content, c.post_id, c.user_id, c.created_at, c.updated_at, users.user_name, users.id FROM comments c
		JOIN users ON c.user_id = users.id
		WHERE c.post_id = $1
		ORDER BY c.created_at 
	`

	rows, err := s.db.QueryContext(ctx, query, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []models.CommentsModel

	for rows.Next() {
		var c models.CommentsModel
		c.User = models.UsersModel{}
		err := rows.Scan(&c.ID, &c.Content, &c.PostID, &c.UserID, &c.CreatedAt, &c.UpdatedAt, &c.User.UserName, &c.User.ID)

		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}
