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
	r.Get("/guitar", h.ListGuitars)
	r.Get("/guitar/{id}", h.GetGuitarByID)
	r.Post("/guitar", h.CreateGuitar)
	r.Put("/guitar/{id}", h.UpdateGuitar)
	r.Patch("/guitar/{id}", h.PatchGuitar)
	r.Delete("/guitar/{id}", h.DeleteGuitar)
}

func (h *GuitarHandler) GetGuitarByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	guitar, err := h.repo.GetGuitarByID(r.Context(), id)
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

func (h *GuitarHandler) ListGuitars(w http.ResponseWriter, r *http.Request) {
	guitars, err := h.repo.ListGuitars(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch guitars", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guitars)
}

func (h *GuitarHandler) CreateGuitar(w http.ResponseWriter, r *http.Request) {
	var guitar Guitar
	err := json.NewDecoder(r.Body).Decode(&guitar)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateGuitar(&guitar); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func (h *GuitarHandler) UpdateGuitar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	var guitar Guitar
	err = json.NewDecoder(r.Body).Decode(&guitar)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateGuitar(&guitar); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	guitar.ID = id

	err = h.repo.UpdateGuitar(r.Context(), &guitar)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Guitar not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update guitar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guitar)
}

func (h *GuitarHandler) PatchGuitar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	guitar, err := h.repo.GetGuitarByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Guitar not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch guitar", http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(guitar)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	guitar.ID = id

	if err := validateGuitar(guitar); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repo.UpdateGuitar(r.Context(), guitar)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Guitar not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update guitar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guitar)
}

func (h *GuitarHandler) DeleteGuitar(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteGuitar(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Guitar not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete guitar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateGuitar(guitar *Guitar) error {
	if strings.TrimSpace(guitar.Manufacturer) == "" {
		return errors.New("Manufacturer is required")
	}
	if guitar.StringCount <= 0 {
		return errors.New("String count must be greater than 0")
	}
	if strings.TrimSpace(guitar.BodyMaterial) == "" {
		return errors.New("Body material is required")
	}
	if strings.TrimSpace(guitar.ManufactureDate) == "" {
		return errors.New("Manufacture date is required")
	}
	return nil
}
