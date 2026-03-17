// Package routes registers all HTTP routes for the application.
// Every new feature's routes must be added here.
package routes

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"my-app/internal/app"
	"my-app/internal/config"
	"my-app/internal/middleware"
)

// New builds and returns the fully configured chi router.
// Add new route groups here when you introduce a new feature.
func New(cfg *config.Config, a *app.App) http.Handler {
	r := chi.NewRouter()

	// ── Global middleware (applied to every request) ──────────────────────────
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS(cfg))
	r.Use(chimiddleware.Timeout(30 * time.Second))

	// ── Swagger UI ─────────────────────────────────────────────────────────────
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// ── Health ─────────────────────────────────────────────────────────────────
	r.Get("/health/live", a.HealthHandler.Live)
	r.Get("/health/ready", a.HealthHandler.Ready)

	// ── API v1 ─────────────────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		// Articles
		r.Post("/articles", a.ArticleHandler.CreateArticle)
		r.Get("/articles/{id}", a.ArticleHandler.GetArticleByID)
		r.Put("/articles/{id}", a.ArticleHandler.UpdateArticle)

		// Add new resource groups here:
		// r.Post("/users",      a.UserHandler.CreateUser)
		// r.Get("/users/{id}",  a.UserHandler.GetUserByID)
	})

	return r
}
