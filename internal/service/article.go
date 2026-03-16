// Package service contains business logic. It sits between handlers and repositories.
// This layer must NOT import pgx, pgxpool, or any sqlc package directly.
package service

import (
	"context"

	"my-app/internal/domain"
	"my-app/internal/repository"
)

// ArticleService owns business logic for the articles resource.
// In this template the logic is intentionally thin — the layer exists so
// future rules (auth checks, rate limits, computed fields) have a home.
type ArticleService struct {
	repo *repository.ArticleRepository
}

// NewArticleService constructs an ArticleService.
func NewArticleService(repo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

// CreateArticle validates business constraints then delegates to the repository.
func (s *ArticleService) CreateArticle(ctx context.Context, req domain.CreateArticleRequest) (domain.ArticleResponse, error) {
	// Business rules go here (e.g. profanity filter, dedup check, quota check).
	return s.repo.CreateArticle(ctx, req)
}

// GetArticleByID delegates retrieval to the repository.
func (s *ArticleService) GetArticleByID(ctx context.Context, id string) (domain.ArticleResponse, error) {
	return s.repo.GetArticleByID(ctx, id)
}
