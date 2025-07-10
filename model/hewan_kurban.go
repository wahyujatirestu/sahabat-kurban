package model

import (
	"time"

	"github.com/google/uuid"
)

type HewanKurban struct {
	ID				uuid.UUID		`db:"id"`
	Jenis			string			`db:"jenis"`
	Berat			float64			`db:"berat"`
	TglPendaftaran	time.Time		`db:"tanggal_pendaftaran"`
	Created_At		time.Time		`db:"created_at"`	
	Updated_At		time.Time		`db:"updated_at"`	
}