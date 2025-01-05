package store

import (
	"context"
	"database/sql"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store/models"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.PostsModel) error
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
