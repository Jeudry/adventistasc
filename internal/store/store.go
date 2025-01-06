package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store/models"
)

var (
	ErrNotFound = errors.New("Resource not found")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.PostsModel) error
		RetrieveById(context.Context, int64) (*models.PostsModel, error)
	}
	Users interface {
		Create(context.Context, *models.UsersModel) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db: db},
		Users: &UsersStore{db: db},
	}
}
