package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
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
	db *sql.DB
}

func NewListingHandler(db *sql.DB) *ListingHandler {
	return &ListingHandler{
		db: db,
	}
}

func (lh *ListingHandler) List(w http.ResponseWriter, r *http.Request) {

	rows, err := lh.db.Query(`
			SELECT id, title, description, price, city, created_at
			FROM listings
			ORDER BY created_at DESC
			LIMIT 100
		`)
	if err != nil {
		log.Printf("query: %v", err)
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
			log.Printf("scan: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		listings = append(listings, l)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(listings); err != nil {
		log.Printf("encode: %v", err)
	}
}

func (lh *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	result, err := lh.db.Exec(
		`DELETE FROM listings WHERE id = $1`,
		id,
	)
	if err != nil {
		log.Printf("delete: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("rows affected: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "listing not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
