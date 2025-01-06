package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store/models"
	"time"
)

var (
	ErrNotFound          = errors.New("Resource not found")
	QueryTimeoutDuration = 5 * time.Second
	ErrConflict          = errors.New("Resource conflict")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.PostsModel) error
		RetrieveById(context.Context, int64) (*models.PostsModel, error)
		Update(context.Context, *models.PostsModel) error
		Delete(context.Context, int64) error
	}
	Users interface {
		Create(context.Context, *models.UsersModel) error
	}
	Comments interface {
		CreatePostComment(context.Context, *models.CommentsModel) error
		RetrieveCommentsByPostId(context.Context, int64) ([]models.CommentsModel, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db: db},
		Users:    &UsersStore{db: db},
		Comments: &CommentsStore{db: db},
	}
}
