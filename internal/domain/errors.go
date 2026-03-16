package domain

import "errors"

// Sentinel errors used across layers. The repository wraps these so the
// handler can distinguish not-found from internal failures without importing pgx.
var (
	ErrNotFound = errors.New("not found")
)
