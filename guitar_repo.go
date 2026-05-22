package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuitarRepo interface {
	GetGuitarByID(ctx context.Context, id uuid.UUID) (*Guitar, error)
	ListGuitars(ctx context.Context) ([]*Guitar, error)
	CreateGuitar(ctx context.Context, guitar *Guitar) error
	UpdateGuitar(ctx context.Context, guitar *Guitar) error
	DeleteGuitar(ctx context.Context, id uuid.UUID) error
}

type guitarRepo struct {
	dbpool *pgxpool.Pool
}

func NewGuitarRepo(dbpool *pgxpool.Pool) GuitarRepo {
	return &guitarRepo{dbpool: dbpool}
}

func (r *guitarRepo) GetGuitarByID(ctx context.Context, id uuid.UUID) (*Guitar, error) {
	query := `SELECT id, manufacturer, string_count, body_material, manufacture_date::text FROM guitar WHERE id = $1`
	guitar := &Guitar{}
	err := r.dbpool.QueryRow(ctx, query, id).Scan(&guitar.ID, &guitar.Manufacturer, &guitar.StringCount, &guitar.BodyMaterial, &guitar.ManufactureDate)
	if err != nil {
		return nil, err
	}
	return guitar, nil
}

func (r *guitarRepo) CreateGuitar(ctx context.Context, guitar *Guitar) error {
	query := `INSERT INTO guitar (id, manufacturer, string_count, body_material, manufacture_date) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.dbpool.Exec(ctx, query, guitar.ID, guitar.Manufacturer, guitar.StringCount, guitar.BodyMaterial, guitar.ManufactureDate)
	return err
}

func (r *guitarRepo) ListGuitars(ctx context.Context) ([]*Guitar, error) {
	query := `SELECT id, manufacturer, string_count, body_material, manufacture_date::text FROM guitar ORDER BY manufacturer`
	rows, err := r.dbpool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	guitars := []*Guitar{}
	for rows.Next() {
		guitar := &Guitar{}
		err := rows.Scan(&guitar.ID, &guitar.Manufacturer, &guitar.StringCount, &guitar.BodyMaterial, &guitar.ManufactureDate)
		if err != nil {
			return nil, err
		}
		guitars = append(guitars, guitar)
	}
	return guitars, rows.Err()
}

func (r *guitarRepo) UpdateGuitar(ctx context.Context, guitar *Guitar) error {
	query := `UPDATE guitar SET manufacturer = $2, string_count = $3, body_material = $4, manufacture_date = $5 WHERE id = $1`
	tag, err := r.dbpool.Exec(ctx, query, guitar.ID, guitar.Manufacturer, guitar.StringCount, guitar.BodyMaterial, guitar.ManufactureDate)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *guitarRepo) DeleteGuitar(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM guitar WHERE id = $1`
	tag, err := r.dbpool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
