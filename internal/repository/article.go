// Package repository owns all database interaction. It wraps sqlc-generated code
// and is the only layer permitted to touch *pgxpool.Pool or sqlc Queries directly.
package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	pkgerrors "github.com/pkg/errors"

	"my-app/internal/db/generated"
	"my-app/internal/domain"
)

// ArticleRepository handles all database operations for articles.
type ArticleRepository struct {
	pool *pgxpool.Pool
}

// NewArticleRepository constructs an ArticleRepository.
func NewArticleRepository(pool *pgxpool.Pool) *ArticleRepository {
	return &ArticleRepository{pool: pool}
}

// CreateArticle inserts a new article inside a transaction and demonstrates
// a multi-step write pattern using pgx transactions with sqlc Queries.
//
// Transaction steps:
//  1. Insert the article row via sqlc-generated CreateArticle query.
//  2. Touch updated_at within the same transaction (simulates a second write).
//  3. Commit. Any error triggers the deferred Rollback automatically.
func (r *ArticleRepository) CreateArticle(ctx context.Context, req domain.CreateArticleRequest) (domain.ArticleResponse, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.ArticleResponse{}, pkgerrors.Wrap(err, "begin tx")
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op after Commit; error irrelevant

	qtx := generated.New(tx)

	article, err := qtx.CreateArticle(ctx, generated.CreateArticleParams{
		Title:  req.Title,
		Body:   req.Body,
		Author: req.Author,
	})
	if err != nil {
		return domain.ArticleResponse{}, pkgerrors.Wrap(err, "insert article")
	}

	// Step 2: second write inside the same transaction — refreshes updated_at.
	// This exists to demonstrate multi-step transaction patterns.
	if _, err := tx.Exec(ctx,
		"UPDATE articles SET updated_at = NOW() WHERE id = $1",
		article.ID,
	); err != nil {
		return domain.ArticleResponse{}, pkgerrors.Wrap(err, "update article timestamp")
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.ArticleResponse{}, pkgerrors.Wrap(err, "commit tx")
	}

	return toArticleResponse(article), nil
}

// GetArticleByID fetches a single article by its UUID string.
// Returns domain.ErrNotFound when no row matches so handlers can return 404.
func (r *ArticleRepository) GetArticleByID(ctx context.Context, id string) (domain.ArticleResponse, error) {
	q := generated.New(r.pool)

	var pgID pgtype.UUID
	if err := pgID.Scan(id); err != nil {
		return domain.ArticleResponse{}, pkgerrors.Wrap(domain.ErrNotFound, "invalid article uuid")
	}

	article, err := q.GetArticleByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ArticleResponse{}, pkgerrors.Wrap(domain.ErrNotFound, "get article by id")
		}
		return domain.ArticleResponse{}, pkgerrors.Wrap(err, "get article by id")
	}

	return toArticleResponse(article), nil
}

// toArticleResponse maps a sqlc Article to the domain response struct.
func toArticleResponse(a generated.Article) domain.ArticleResponse {
	return domain.ArticleResponse{
		ID:        uuid.UUID(a.ID.Bytes).String(),
		Title:     a.Title,
		Body:      a.Body,
		Author:    a.Author,
		CreatedAt: a.CreatedAt.Time,
		UpdatedAt: a.UpdatedAt.Time,
	}
}
