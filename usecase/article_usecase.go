package usecase

import (
	"context"
	"errors"
	"github.com/jpdel518/go-onionarch/domain/model"
	"github.com/jpdel518/go-onionarch/domain/repository"
	"github.com/jpdel518/go-onionarch/domain/service"
	"log"
	"sync"
	"time"
)

type ArticleUsecase interface {
	Fetch(ctx context.Context, num int64) ([]model.Article, error)
	GetByID(ctx context.Context, id int64) (model.Article, error)
	Update(ctx context.Context, ar *model.Article) error
	GetByTitle(ctx context.Context, title string) (model.Article, error)
	Store(context.Context, *model.Article) error
	Delete(ctx context.Context, id int64) error
}

type articleUsecase struct {
	articleRepo    repository.ArticleRepository
	authorRepo     repository.AuthorRepository
	articleService service.ArticleService
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object
func NewArticleUsecase(a repository.ArticleRepository, ar repository.AuthorRepository, as service.ArticleService, timeout time.Duration) ArticleUsecase {
	return &articleUsecase{
		articleRepo:    a,
		authorRepo:     ar,
		articleService: as,
		contextTimeout: timeout,
	}
}

// fillAuthorDetails will fill up author details that is concerned with a parameter of article object
func (a *articleUsecase) fillAuthorDetails(c context.Context, data []model.Article) ([]model.Article, error) {
	// TODO errgroupを使ってエラーをハンドリング + contextで後続の処理をキャンセルするようにした方がいいかも
	var wg sync.WaitGroup

	// Get the author's id
	mapAuthors := map[int64]model.Author{}

	for _, article := range data {
		mapAuthors[article.Author.ID] = model.Author{}
	}
	// Using goroutine to fetch the author's detail
	chanAuthor := make(chan model.Author)

	for authorID := range mapAuthors {
		wg.Add(1)
		go func(authorID int64) {
			defer wg.Done()
			res, err := a.authorRepo.GetByID(c, authorID)
			chanAuthor <- res
			if err != nil {
				log.Printf("failed author %d, error caused: %+v", authorID, err)
			}
		}(authorID)
	}

	// release channel
	go func() {
		wg.Wait()
		close(chanAuthor)
	}()

	// retrieve
	for author := range chanAuthor {
		if author != (model.Author{}) {
			mapAuthors[author.ID] = author
		}
	}
	wg.Wait()

	// merge the author's data
	for index, item := range data {
		if a, ok := mapAuthors[item.Author.ID]; ok {
			data[index].Author = a
		}
	}
	return data, nil
}

// Fetch will retrieve articles
func (a *articleUsecase) Fetch(c context.Context, num int64) ([]model.Article, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err := a.articleRepo.Fetch(ctx, num)
	if err != nil {
		return nil, err
	}

	res, err = a.fillAuthorDetails(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetByID will find an article by id
func (a *articleUsecase) GetByID(c context.Context, id int64) (model.Article, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return model.Article{}, err
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return model.Article{}, err
	}
	res.Author = resAuthor
	return res, nil
}

// Update will update an article
func (a *articleUsecase) Update(c context.Context, ar *model.Article) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

// GetByTitle will find article by title
func (a *articleUsecase) GetByTitle(c context.Context, title string) (model.Article, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err := a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return model.Article{}, err
	}

	resAuthor, err := a.authorRepo.GetByID(ctx, res.Author.ID)
	if err != nil {
		return res, err
	}

	res.Author = resAuthor
	return res, nil
}

// Store will register an article
func (a *articleUsecase) Store(c context.Context, m *model.Article) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, err := a.GetByTitle(ctx, m.Title)
	if existedArticle != (model.Article{}) {
		return err
	}

	err = a.articleRepo.Store(ctx, m)
	if err != nil {
		return err
	}
	// domain logic cvb
	a.articleService.ReleaseRecommend(ctx, m.ID)
	return err
}

// Delete will delete an article by id
func (a *articleUsecase) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existedArticle == (model.Article{}) {
		return errors.New("not found")
	}
	return a.articleRepo.Delete(ctx, id)
}
