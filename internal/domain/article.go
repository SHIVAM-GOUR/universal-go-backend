package domain

import "time"

// CreateArticleRequest is the validated request body for creating an article.
type CreateArticleRequest struct {
	Title  string `json:"title"  validate:"required,min=3,max=255"`
	Body   string `json:"body"   validate:"required,min=10"`
	Author string `json:"author" validate:"required,min=2,max=100"`
}

// ArticleResponse is the JSON representation of an article returned to clients.
type ArticleResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ErrorResponse is the standard error envelope returned on all failed requests.
type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
	Code  string `json:"code"  example:"INTERNAL_ERROR"`
}

type UpdateArticleRequest struct {
	Title string `json:"title" validate:"required,min=3,max=255"`
	Body  string `json:"body" validate:"required,min=10"`
}
