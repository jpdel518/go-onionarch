package repository

import (
	"context"
	"github.com/jpdel518/go-onionarch/domain/model"
)

type ArticleRepository interface {
	Fetch(ctx context.Context, num int64) (res []model.Article, err error)
	GetByID(ctx context.Context, id int64) (model.Article, error)
	GetByTitle(ctx context.Context, title string) (model.Article, error)
	Update(ctx context.Context, ar *model.Article) error
	Store(ctx context.Context, a *model.Article) error
	Delete(ctx context.Context, id int64) error
}
