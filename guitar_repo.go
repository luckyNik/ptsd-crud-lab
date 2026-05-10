package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GuitarRepo interface {
	GetGuitarByID(ctx context.Context, id string) (*Guitar, error)
	CreateGuitar(ctx context.Context, guitar *Guitar) error
}

type guitarRepo struct {
	dbpool *pgxpool.Pool
}

func NewGuitarRepo(dbpool *pgxpool.Pool) GuitarRepo {
	return &guitarRepo{dbpool: dbpool}
}

func (r *guitarRepo) GetGuitarByID(ctx context.Context, id string) (*Guitar, error) {
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
