package service

import (
	"context"
	"github.com/jpdel518/go-onionarch/domain/repository"
	"log"
	"time"
)

type ArticleService interface {
	ReleaseRecommend(ctx context.Context, id int64)
}

type articleServiceImpl struct {
	repository.ArticleRepository
}

// NewArticleService create article domain service instance
func NewArticleService(r repository.ArticleRepository) ArticleService {
	return &articleServiceImpl{r}
}

// ReleaseRecommend will register recommend the article when the article is stored
func (s *articleServiceImpl) ReleaseRecommend(ctx context.Context, id int64) {
	article, err := s.ArticleRepository.GetByID(ctx, id)
	if err != nil {
		log.Printf("failed get recommend article: %s", err)
	}
	article.Recommend = "It's new article"
	article.RecommendedAt = time.Now()
	err = s.Update(ctx, &article)
	if err != nil {
		log.Printf("failed update recomend article: %s", err)
	}
}
