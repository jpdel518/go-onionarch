package rdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jpdel518/go-onionarch/domain/model"
	"github.com/jpdel518/go-onionarch/domain/repository"
	"log"
)

type articleRepository struct {
	Conn *sql.DB
}

// NewArticleRepository will create an object that represent the article.Repository interface
func NewArticleRepository(conn *sql.DB) repository.ArticleRepository {
	return &articleRepository{conn}
}

// fetch is common fetch logic
func (m *articleRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]model.Article, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Println(errRow)
		}
	}()

	result := make([]model.Article, 0)
	for rows.Next() {
		t := model.Article{}
		authorID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		t.Author = model.Author{
			ID: authorID,
		}
		result = append(result, t)
	}

	return result, nil
}

// Fetch will get number of articles
func (m *articleRepository) Fetch(ctx context.Context, num int64) ([]model.Article, error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
  						FROM article WHERE created_at > ? ORDER BY created_at LIMIT ? `

	res, err := m.fetch(ctx, query, num)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetByID will get an article by article id
func (m *articleRepository) GetByID(ctx context.Context, id int64) (model.Article, error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
  						FROM article WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return model.Article{}, err
	}

	res := model.Article{}
	if len(list) > 0 {
		res = list[0]
	} else {
		return model.Article{}, errors.New(fmt.Sprintf("article not found by ID: %d", id))
	}

	return res, nil
}

// GetByTitle will get an article by title
func (m *articleRepository) GetByTitle(ctx context.Context, title string) (model.Article, error) {
	query := `SELECT id,title,content, author_id, updated_at, created_at
  						FROM article WHERE title = ?`

	res := model.Article{}
	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return res, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, errors.New(fmt.Sprintf("Article not found by Title: %s", title))
	}
	return res, nil
}

// Store will register an article
func (m *articleRepository) Store(ctx context.Context, a *model.Article) error {
	query := `INSERT  article SET title=? , content=? , author_id=?, updated_at=? , created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = lastID
	return nil
}

// Delete will delete an article by specified id
func (m *articleRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM article WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return err
	}

	return nil
}

// Update will update an article by specified id
func (m *articleRepository) Update(ctx context.Context, ar *model.Article) error {
	query := `UPDATE article set title=?, content=?, author_id=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return err
	}

	return nil
}
