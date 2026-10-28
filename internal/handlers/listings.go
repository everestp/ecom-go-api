package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"

	"log/slog"
	"net/http"

	"github.com/everestp/ecom-go-api/internal/httpx"
	"github.com/everestp/ecom-go-api/internal/middleware"
)

type listing struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	City        string `json:"city"`
	CreatedAt   string `json:"created_at"`
}

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

func (lh *ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	//request scoped context
	ctx := r.Context()

	rows, err := lh.db.QueryContext(ctx, `
			SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100
		`)
	if err != nil {
		// lh.logger.Error("query: %v", err)
		lh.logger.Error("query: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	listings := []listing{}

	for rows.Next() {
		var l listing

		if err := rows.Scan(
			&l.ID,
			&l.Title,
			&l.Description,
			&l.Price,
			&l.City,
			&l.CreatedAt,
		); err != nil {
			lh.logger.Error("scan: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		lh.logger.Error("rows: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(listings); err != nil {
		lh.logger.Error("encode: %v", err)
	}
}

func (lh *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// requestId := ctx.Value("requestCtxId").(string)
	requestId := middleware.RequestIDFromContext(ctx)
	id := r.PathValue("id")

	result, err := lh.db.ExecContext(ctx,
		`DELETE FROM listing WHERE id = $1`,
		id,
	)
	if err != nil {
		lh.logger.Error("delete failed", "listing_id", id, "request_id", requestId, "err", err)

		httpx.Error(w, http.StatusInternalServerError, "Something went wring", "internal_error")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		lh.logger.Error("rows affected: %v", err)
		// http.Error(w, "internal error", http.StatusInternalServerError)
		httpx.Error(w, http.StatusInternalServerError, "Something went wring", httpx.CodeInternalError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "listing not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (lh *ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// requestId := ctx.Value("requestCtxId").(string)
	requestId := middleware.RequestIDFromContext(ctx)

	var req CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.Error("failed to decode", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	if err := req.Validate(); err != nil {
		var verr ValidationError
		errors.As(err, &verr)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.CodeValidationFailed, verr.Field)
		return
	}

	row := lh.db.QueryRowContext(ctx, `INSERT INTO listings (title , description , price ,city) VALUES ($1,$2,$3,$4) RETURNING id ,title ,created_at`, req.Title, req.Description, req.Price, req.City)
	// var id string
	var out CreateListingResponse
	if err := row.Scan(&out.ID, &out.Title, &out.CreatedAt); err != nil {
		lh.logger.Error("failed to insert", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	lh.logger.Info("listing created", "request_id", requestId, "listing_id", out.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		lh.logger.Error("encode: %v", err)
	}

}
