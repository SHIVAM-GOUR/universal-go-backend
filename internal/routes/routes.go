// Package routes registers all HTTP routes for the application.
// Every new feature's routes must be added here.
package routes

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"my-app/internal/config"
	"my-app/internal/handler"
	"my-app/internal/middleware"
)

// New builds and returns the fully configured chi router.
func New(cfg *config.Config, healthHandler *handler.HealthHandler, articleHandler *handler.ArticleHandler) http.Handler {
	r := chi.NewRouter()

	// ── Global middleware (applied to every request) ──────────────────────────
	r.Use(chimiddleware.RequestID)   // injects X-Request-ID header
	r.Use(chimiddleware.RealIP)      // reads X-Forwarded-For / X-Real-IP
	r.Use(chimiddleware.Recoverer)   // catches panics, returns 500
	r.Use(middleware.CORS(cfg))      // cross-origin headers
	r.Use(chimiddleware.Timeout(30 * time.Second)) // per-request deadline

	// ── Swagger UI ─────────────────────────────────────────────────────────────
	// Access at: GET /swagger/index.html
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// ── Health endpoints ───────────────────────────────────────────────────────
	r.Get("/health/live", healthHandler.Live)
	r.Get("/health/ready", healthHandler.Ready)

	// ── API v1 ─────────────────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		// Articles
		r.Post("/articles", articleHandler.CreateArticle)
		r.Get("/articles/{id}", articleHandler.GetArticleByID)
	})

	return r
}
