package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"my-app/internal/domain"
	"my-app/internal/service"
)

// ArticleHandler handles HTTP requests for the articles resource.
type ArticleHandler struct {
	svc *service.ArticleService
}

// NewArticleHandler constructs an ArticleHandler.
func NewArticleHandler(svc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// CreateArticle godoc
//
//	@Summary		Create an article
//	@Description	Creates a new article using a DB transaction
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			request	body		domain.CreateArticleRequest	true	"Article payload"
//	@Success		201		{object}	domain.ArticleResponse
//	@Failure		400		{object}	domain.ErrorResponse
//	@Failure		500		{object}	domain.ErrorResponse
//	@Router			/api/v1/articles [post]
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", "INVALID_BODY")
		return
	}

	if err := validate.Struct(req); err != nil {
		writeValidationError(w, err)
		return
	}

	resp, err := h.svc.CreateArticle(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create article", "INTERNAL_ERROR")
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// GetArticleByID godoc
//
//	@Summary		Get an article
//	@Description	Retrieves a single article by UUID
//	@Tags			articles
//	@Produce		json
//	@Param			id	path		string	true	"Article UUID"
//	@Success		200	{object}	domain.ArticleResponse
//	@Failure		404	{object}	domain.ErrorResponse
//	@Failure		500	{object}	domain.ErrorResponse
//	@Router			/api/v1/articles/{id} [get]
func (h *ArticleHandler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing article id", "MISSING_ID")
		return
	}

	resp, err := h.svc.GetArticleByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "article not found", "NOT_FOUND")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to retrieve article", "INTERNAL_ERROR")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// UpdateArticle godoc
//
//	@Summary                Update an article
//	@Tags                   articles
//	@Accept                 json
//	@Produce                json
//	@Param                  id              path            string                                          true    "Article UUID"
//	@Param                  request body            domain.UpdateArticleRequest     true    "Update payload"
//	@Success                200             {object}        domain.ArticleResponse
//	@Failure                400             {object}        domain.ErrorResponse
//	@Failure                404             {object}        domain.ErrorResponse
//	@Router                 /api/v1/articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req domain.UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", "INVALID BODY")
		return
	}

	if err := validate.Struct(req); err != nil {
		writeValidationError(w, err)
		return
	}

	resp, err := h.svc.UpdateArticle(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "article not found", "NOT FOUND")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update article", "INTERNAL_ERROR")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
