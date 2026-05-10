package main

import (
	"github.com/google/uuid"
)

type Guitar struct {
	ID              uuid.UUID `json:"id"`
	Manufacturer    string    `json:"manufacturer"`
	StringCount     int       `json:"string_count"`
	BodyMaterial    string    `json:"body_material"`
	ManufactureDate string    `json:"manufacture_date"`
}
