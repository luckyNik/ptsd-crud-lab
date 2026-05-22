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

func (m *mockGuitarRepo) ListGuitars(ctx context.Context) ([]*Guitar, error) {
	guitars := []*Guitar{}
	for _, g := range m.guitars {
		guitars = append(guitars, g)
	}
	return guitars, nil
}

func (m *mockGuitarRepo) UpdateGuitar(ctx context.Context, guitar *Guitar) error {
	if _, exists := m.guitars[guitar.ID]; !exists {
		return pgx.ErrNoRows
	}
	m.guitars[guitar.ID] = guitar
	return nil
}

func (m *mockGuitarRepo) DeleteGuitar(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.guitars[id]; !exists {
		return pgx.ErrNoRows
	}
	delete(m.guitars, id)
	return nil
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

func TestListHandler_Success(t *testing.T) {
	repo := newMockGuitarRepo()
	id := uuid.New()
	repo.guitars[id] = &Guitar{
		ID:              id,
		Manufacturer:    "Gibson",
		StringCount:     6,
		BodyMaterial:    "Mahogany",
		ManufactureDate: "2023-06-15",
	}

	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/guitar", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got []Guitar
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("Expected 1 guitar, got %d", len(got))
	}
	if got[0].ID != id {
		t.Fatalf("Expected id %s, got %s", id, got[0].ID)
	}
}

func TestListHandler_Empty(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/guitar", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got []Guitar
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Expected 0 guitars, got %d", len(got))
	}
}

func TestUpdateHandler_Success(t *testing.T) {
	repo := newMockGuitarRepo()
	id := uuid.New()
	repo.guitars[id] = &Guitar{
		ID:              id,
		Manufacturer:    "Fender",
		StringCount:     6,
		BodyMaterial:    "Alder",
		ManufactureDate: "2020-01-01",
	}

	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	updated := Guitar{
		Manufacturer:    "Fender",
		StringCount:     7,
		BodyMaterial:    "Ash",
		ManufactureDate: "2024-02-02",
	}
	body, _ := json.Marshal(updated)
	req := httptest.NewRequest(http.MethodPut, "/guitar/"+id.String(), bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if repo.guitars[id].StringCount != 7 || repo.guitars[id].BodyMaterial != "Ash" {
		t.Fatalf("Guitar not updated: %+v", repo.guitars[id])
	}
}

func TestUpdateHandler_InvalidUUID(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	body, _ := json.Marshal(&Guitar{Manufacturer: "X", StringCount: 6, BodyMaterial: "Y", ManufactureDate: "2024-01-01"})
	req := httptest.NewRequest(http.MethodPut, "/guitar/not-a-uuid", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateHandler_NotFound(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	body, _ := json.Marshal(&Guitar{Manufacturer: "X", StringCount: 6, BodyMaterial: "Y", ManufactureDate: "2024-01-01"})
	req := httptest.NewRequest(http.MethodPut, "/guitar/"+uuid.New().String(), bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUpdateHandler_InvalidBody(t *testing.T) {
	repo := newMockGuitarRepo()
	id := uuid.New()
	repo.guitars[id] = &Guitar{ID: id, Manufacturer: "Fender", StringCount: 6, BodyMaterial: "Alder", ManufactureDate: "2020-01-01"}

	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	body, _ := json.Marshal(&Guitar{Manufacturer: "", StringCount: 0, BodyMaterial: "", ManufactureDate: ""})
	req := httptest.NewRequest(http.MethodPut, "/guitar/"+id.String(), bytes.NewReader(body))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestPatchHandler_Success(t *testing.T) {
	repo := newMockGuitarRepo()
	id := uuid.New()
	repo.guitars[id] = &Guitar{
		ID:              id,
		Manufacturer:    "Fender",
		StringCount:     6,
		BodyMaterial:    "Alder",
		ManufactureDate: "2020-01-01",
	}

	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPatch, "/guitar/"+id.String(), strings.NewReader(`{"string_count": 12}`))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if repo.guitars[id].StringCount != 12 {
		t.Fatalf("Expected string_count 12, got %d", repo.guitars[id].StringCount)
	}
	if repo.guitars[id].Manufacturer != "Fender" {
		t.Fatalf("Expected manufacturer unchanged, got %s", repo.guitars[id].Manufacturer)
	}
}

func TestPatchHandler_NotFound(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPatch, "/guitar/"+uuid.New().String(), strings.NewReader(`{"string_count": 7}`))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestPatchHandler_InvalidUUID(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPatch, "/guitar/not-a-uuid", strings.NewReader(`{"string_count": 7}`))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestPatchHandler_InvalidBody(t *testing.T) {
	repo := newMockGuitarRepo()
	id := uuid.New()
	repo.guitars[id] = &Guitar{ID: id, Manufacturer: "Fender", StringCount: 6, BodyMaterial: "Alder", ManufactureDate: "2020-01-01"}

	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodPatch, "/guitar/"+id.String(), strings.NewReader(`{"manufacturer": ""}`))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestDeleteHandler_Success(t *testing.T) {
	repo := newMockGuitarRepo()
	id := uuid.New()
	repo.guitars[id] = &Guitar{ID: id, Manufacturer: "Fender", StringCount: 6, BodyMaterial: "Alder", ManufactureDate: "2020-01-01"}

	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodDelete, "/guitar/"+id.String(), nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	if _, exists := repo.guitars[id]; exists {
		t.Fatalf("Guitar was not deleted")
	}
}

func TestDeleteHandler_NotFound(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodDelete, "/guitar/"+uuid.New().String(), nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestDeleteHandler_InvalidUUID(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodDelete, "/guitar/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
