package rdb

import (
	"context"
	"database/sql"
	"github.com/jpdel518/go-onionarch/domain/model"
	"github.com/jpdel518/go-onionarch/domain/repository"
)

type authorRepo struct {
	DB *sql.DB
}

// NewAuthorRepository will create an implementation of author.Repository
func NewAuthorRepository(db *sql.DB) repository.AuthorRepository {
	return &authorRepo{
		DB: db,
	}
}

// getOne
func (m *authorRepo) getOne(ctx context.Context, query string, args ...interface{}) (model.Author, error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return model.Author{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	res := model.Author{}

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return model.Author{}, err
	}
	return res, nil
}

func (m *authorRepo) GetByID(ctx context.Context, id int64) (model.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id=?`
	return m.getOne(ctx, query, id)
}
