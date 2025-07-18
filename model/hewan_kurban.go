package model

import (
	"time"

	"github.com/google/uuid"
)

type JenisHewan string

const (
	Sapi    JenisHewan = "sapi"
	Kambing JenisHewan = "kambing"
	Domba   JenisHewan = "domba"
)

type HewanKurban struct {
	ID                 	uuid.UUID   `db:"id"`
	Jenis              	JenisHewan  `db:"jenis"`
	Berat              	float64     `db:"berat"`
	Harga              	float64     `db:"harga"`
	IsPrivate          	bool        `db:"is_private"`
	TanggalPendaftaran 	time.Time   `db:"tanggal_pendaftaran"`
	Created_At          time.Time   `db:"created_at"`
	Updated_At          time.Time   `db:"updated_at"`
}
