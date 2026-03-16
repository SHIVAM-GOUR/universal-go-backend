// Package app wires every layer of the application together.
// This is the only place where repositories, services, and handlers are
// constructed and connected. main.go never needs to change when you add a
// new feature — only this file and routes/routes.go do.
package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"my-app/internal/config"
	"my-app/internal/handler"
	"my-app/internal/repository"
	"my-app/internal/service"
)

// App holds every HTTP handler in the application.
// Add a new exported field here whenever you introduce a new feature.
type App struct {
	HealthHandler  *handler.HealthHandler
	ArticleHandler *handler.ArticleHandler
	// UserHandler *handler.UserHandler  ← add future handlers here
}

// New wires all layers (repository → service → handler) and returns a
// fully initialised App ready to be handed to the router.
//
// To add a new feature:
//  1. Create repository, service, handler types in their respective packages.
//  2. Instantiate them here following the same pattern as articles.
//  3. Assign the handler to a new field on App.
//  4. Register its routes in internal/routes/routes.go.
func New(cfg *config.Config, pool *pgxpool.Pool) *App {
	// ── Articles ───────────────────────────────────────────────────────────────
	articleRepo := repository.NewArticleRepository(pool)
	articleSvc := service.NewArticleService(articleRepo)

	// ── Future features go below this line ────────────────────────────────────
	// userRepo := repository.NewUserRepository(pool)
	// userSvc  := service.NewUserService(userRepo)

	return &App{
		HealthHandler:  handler.NewHealthHandler(cfg, pool),
		ArticleHandler: handler.NewArticleHandler(articleSvc),
		// UserHandler: handler.NewUserHandler(userSvc),
	}
}
