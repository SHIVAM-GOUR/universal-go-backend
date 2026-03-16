package handler

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"my-app/internal/config"
	"my-app/internal/domain"
)

// HealthHandler serves liveness and readiness probes.
type HealthHandler struct {
	cfg  *config.Config
	pool *pgxpool.Pool
}

// NewHealthHandler constructs a HealthHandler.
func NewHealthHandler(cfg *config.Config, pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{cfg: cfg, pool: pool}
}

// Live godoc
//
//	@Summary		Liveness check
//	@Description	Returns 200 if the server process is running
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	domain.HealthResponse
//	@Router			/health/live [get]
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, domain.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Env:       h.cfg.AppEnv,
		Version:   "1.0.0",
	})
}

// Ready godoc
//
//	@Summary		Readiness check
//	@Description	Returns 200 if server and DB are both reachable
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	domain.HealthResponse
//	@Failure		503	{object}	domain.ErrorResponse
//	@Router			/health/ready [get]
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if err := h.pool.Ping(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database not reachable", "DB_UNAVAILABLE")
		return
	}

	writeJSON(w, http.StatusOK, domain.HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Env:       h.cfg.AppEnv,
		Version:   "1.0.0",
	})
}
