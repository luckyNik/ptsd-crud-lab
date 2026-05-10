// TODO: implement other CRUD endpoints according to REST (GET list, PUT, PATCH, DELETE)
package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GuitarHandler struct {
	repo GuitarRepo
}

func NewGuitarHandler(repo GuitarRepo) *GuitarHandler {
	return &GuitarHandler{repo: repo}
}

func (h *GuitarHandler) RegisterRoutes(r chi.Router) {
	r.Get("/guitar/{id}", h.GetGuitarByID)
	r.Post("/guitar", h.CreateGuitar)
}

func (h *GuitarHandler) GetGuitarByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	guitar, err := h.repo.GetGuitarByID(r.Context(), id.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Guitar not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to fetch guitar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guitar)
}

func (h *GuitarHandler) CreateGuitar(w http.ResponseWriter, r *http.Request) {
	var guitar Guitar
	err := json.NewDecoder(r.Body).Decode(&guitar)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(guitar.Manufacturer) == "" {
		http.Error(w, "Manufacturer is required", http.StatusBadRequest)
		return
	}

	if guitar.StringCount <= 0 {
		http.Error(w, "String count must be greater than 0", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(guitar.BodyMaterial) == "" {
		http.Error(w, "Body material is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(guitar.ManufactureDate) == "" {
		http.Error(w, "Manufacture date is required", http.StatusBadRequest)
		return
	}

	guitar.ID = uuid.New()

	err = h.repo.CreateGuitar(r.Context(), &guitar)
	if err != nil {
		http.Error(w, "Failed to create guitar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(guitar)
}
