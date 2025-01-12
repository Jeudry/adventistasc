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

func (s *PostsStore) GetUserFeed(ctx context.Context, userId int64, fq models.PaginatedFeedQueryModel) ([]models.PostWithMetadata, error) {
	var query string

	baseQuery := `SELECT p.id, p.user_id, p.title, p.content, p.created_at, 
            p.updated_at, p.version, p.tags, u.user_name, 
            u.email, p.user_id, COUNT(c.id) AS comments_count 
            FROM posts p 
            LEFT JOIN comments c ON p.id = c.post_id
            LEFT JOIN users u on p.user_id = u.id
            JOIN followers f ON f.follower_id = p.user_id
            WHERE f.user_id = $1`

	if fq.Search != "" {
		query = baseQuery + `
        AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
    `
	} else {
		query = baseQuery + ` 
        OR p.user_id = $1
    `
	}

	if len(fq.Tags) > 0 {
		if fq.Search != "" {
			query += `AND (p.tags @> $4 OR $4 = '{}')`
		} else {
			query += `AND (p.tags @> $5 OR $5 = '{}')`
		}
	}

	query += `
    	GROUP BY p.id, u.user_name, u.email, u.id
    	ORDER BY p.created_at ` + fq.Sort + `
    	LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var rows *sql.Rows
	var err error

	if fq.Search != "" {
		rows, err = s.db.QueryContext(ctx, query, userId, fq.Limit, fq.Offset, fq.Search)
	} else if len(fq.Tags) > 0 {
		rows, err = s.db.QueryContext(ctx, query, userId, fq.Limit, fq.Offset, pq.Array(fq.Tags))
	} else {
		rows, err = s.db.QueryContext(ctx, query, userId, fq.Limit, fq.Offset)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []models.PostWithMetadata

	for rows.Next() {
		var p models.PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.UserName,
			&p.User.Email,
			&p.User.ID,
			&p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil
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

	ctx, timeout := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer timeout()

	err := s.db.QueryRowContext(ctx, query, postUpdated.Content, postUpdated.Title, postUpdated.Version).Scan(
		&postUpdated.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrConflict
		default:
			return err
		}
	}

	return err
}

func (s *PostsStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, timeout := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer timeout()

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
