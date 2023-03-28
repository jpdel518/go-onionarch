package repository

import (
	"context"
	"github.com/jpdel518/go-onionarch/domain/model"
)

type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (model.Author, error)
}
