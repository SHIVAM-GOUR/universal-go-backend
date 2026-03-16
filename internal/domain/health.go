// Package domain contains all request/response structs and domain-level sentinel errors.
// No business logic lives here.
package domain

// HealthResponse is returned by liveness and readiness endpoints.
type HealthResponse struct {
	Status    string `json:"status"    example:"ok"`
	Timestamp string `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Env       string `json:"env"       example:"development"`
	Version   string `json:"version"   example:"1.0.0"`
}
