package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type mockGuitarRepo struct {
	guitars map[uuid.UUID]*Guitar
}

func newMockGuitarRepo() *mockGuitarRepo {
	return &mockGuitarRepo{
		guitars: make(map[uuid.UUID]*Guitar),
	}
}

func (m *mockGuitarRepo) CreateGuitar(ctx context.Context, guitar *Guitar) error {
	m.guitars[guitar.ID] = guitar
	return nil
}

func (m *mockGuitarRepo) GetGuitarByID(ctx context.Context, id uuid.UUID) (*Guitar, error) {
	guitar, exists := m.guitars[id]
	if !exists {
		return nil, pgx.ErrNoRows
	}
	return guitar, nil
}

func TestCreatrHandler_InvalidJSON(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/guitar", strings.NewReader("invalid json"))
	rec := httptest.NewRecorder()
	handler.CreateGuitar(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateHandler_EmptyManufacturer(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)
	guitar := &Guitar{
		ID:              uuid.New(),
		Manufacturer:    "",
		StringCount:     6,
		BodyMaterial:    "Wood",
		ManufactureDate: "2024-01-01",
	}

	body, _ := json.Marshal(guitar)
	req := httptest.NewRequest(http.MethodPost, "/guitar", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.CreateGuitar(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetHandler_Success(t *testing.T) {
	repo := newMockGuitarRepo()
	id := uuid.New()
	repo.guitars[id] = &Guitar{
		ID:              id,
		Manufacturer:    "Ibanez",
		StringCount:     6,
		BodyMaterial:    "Oak",
		ManufactureDate: "2024-01-01",
	}

	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/guitar/"+id.String(), nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got Guitar
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if got.ID != id {
		t.Fatalf("Expected id %s, got %s", id, got.ID)
	}
}

func TestGetHandler_NotFound(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/guitar/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
