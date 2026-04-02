package handler

import (
	"net/http"
	"strconv"

	"github.com/oack-io/poke-store/internal/data"
	"github.com/oack-io/poke-store/internal/model"
)

type Store struct {
	catalog *data.Catalog
}

func NewStore(catalog *data.Catalog) *Store {
	return &Store{catalog: catalog}
}

func (h *Store) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	typeFilter := r.URL.Query().Get("type")

	results := h.catalog.Search(query)

	if typeFilter != "" {
		var filtered []model.Pokemon
		for _, p := range results {
			for _, t := range p.Type {
				if t == typeFilter {
					filtered = append(filtered, p)
					break
				}
			}
		}
		results = filtered
	}

	// Catalog is static — let CDN and browsers cache it.
	w.Header().Set("Cache-Control", "public, max-age=60, stale-while-revalidate=300")

	JSON(w, http.StatusOK, map[string]any{
		"pokemon": results,
		"count":   len(results),
	})
}

func (h *Store) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid pokemon id")
		return
	}

	pokemon, ok := h.catalog.Get(id)
	if !ok {
		Error(w, http.StatusNotFound, "pokemon not found")
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=60, stale-while-revalidate=300")
	JSON(w, http.StatusOK, pokemon)
}
