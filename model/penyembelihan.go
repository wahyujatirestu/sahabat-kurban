package model

import (
	"time"

	"github.com/google/uuid"
)

type Penyembelihan struct {
	ID 					uuid.UUID		`db:"id"`
	HewanID				uuid.UUID		`db:"hewan_id"`
	JenisHewan			JenisHewan		`db:"jenis_hewan"`
	TglPenyembelihan 	time.Time		`db:"tanggal_penyembelihan"`
	Lokasi				string			`db:"lokasi"`
	UrutanRencana		int				`db:"urutan_rencana"`
	UrutanAktual		*int 			`db:"Urutan_aktual"`
	Created_At			time.Time		`db:"created_at"`
	Updated_At			time.Time		`db:"updated_at"`
}