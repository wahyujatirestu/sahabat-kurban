package model

import (
	"time"

	"github.com/google/uuid"
)

type DistribusiDaging struct {
	ID                	uuid.UUID		`db:"id"`
	PenerimaID        	uuid.UUID		`db:"penerima_id"`
	HewanID           	uuid.UUID		`db:"hewan_id"`
	JumlahPaket       	int				`db:"jumlah_paket"`
	TanggalDistribusi 	time.Time		`db:"tanggal_distribusi"`
	Created_At         	time.Time		`db:"created_at"`
	Updated_At         	time.Time		`db:"updated_at"`
}
