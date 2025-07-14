package model

import (
	"time"

	"github.com/google/uuid"
)

type DistribusiDaging struct {
	ID                uuid.UUID
	PenerimaID        uuid.UUID
	HewanID           uuid.UUID
	JumlahPaket       int
	TanggalDistribusi time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
