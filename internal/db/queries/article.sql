-- name: CreateArticle :one
INSERT INTO articles (title, body, author)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetArticleByID :one
SELECT * FROM articles
WHERE id = $1;

-- name: ListArticles :many
SELECT * FROM articles
ORDER BY created_at DESC;

-- name: UpdateArticle :one
UPDATE articles
SET title = $2, body = $3, updated_at = now()
WHERE id = $1
RETURNING *;