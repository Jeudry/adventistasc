package store

import (
	"context"
	"database/sql"
)

type FollowersStore struct {
	db *sql.DB
}

func (s *FollowersStore) ToggleUserFollow(ctx context.Context, followerId int64, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	deleteQuery := `
        DELETE FROM followers 
        WHERE follower_id = $1 AND user_id = $2
        RETURNING *;
    `

	result, err := s.db.ExecContext(ctx, deleteQuery, followerId, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		insertQuery := `
            INSERT INTO followers (follower_id, user_id)
            VALUES ($1, $2);
        `
		_, err = s.db.ExecContext(ctx, insertQuery, followerId, userId)
		if err != nil {
			return err
		}
	}

	return nil
}
