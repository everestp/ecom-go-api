package handlers

import (
	"database/sql"
	"encoding/json"

	"log/slog"
	"net/http"

	"github.com/everestp/ecom-go-api/internal/middleware"
)

type listing struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`
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
		`DELETE FROM listings WHERE id = $1`,
		id,
	)
	if err != nil {
		lh.logger.Error("delete failed", "listing_id", id, "request_id", requestId, "err", err)

		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		lh.logger.Error("rows affected: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "listing not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
