package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"my-app/internal/domain"
)

// validate is a package-level singleton. Creating a validator is expensive;
// reuse it across all handlers.
var validate = validator.New()

// writeJSON serialises v as JSON and writes it with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a standard error envelope JSON response.
func writeError(w http.ResponseWriter, status int, message, code string) {
	writeJSON(w, status, domain.ErrorResponse{
		Error: message,
		Code:  code,
	})
}

// writeValidationError formats validator.ValidationErrors into a readable
// error message and writes a 400 response.
func writeValidationError(w http.ResponseWriter, err error) {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		writeError(w, http.StatusBadRequest, "invalid request", "VALIDATION_ERROR")
		return
	}

	msgs := make([]string, 0, len(ve))
	for _, e := range ve {
		msgs = append(msgs, e.Field()+": failed '"+e.Tag()+"' validation")
	}

	writeError(w, http.StatusBadRequest, strings.Join(msgs, "; "), "VALIDATION_ERROR")
}
