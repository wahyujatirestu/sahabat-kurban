package model

import "github.com/google/uuid"

type TotalPembayaranPerHewan struct {
	HewanID     uuid.UUID
	Jenis       string
	HargaTarget float64
	TotalMasuk  float64
}

type ProgressPembayaran struct {
	PekurbanID   uuid.UUID
	NamaPekurban string
	PorsiTotal   float64
	TotalTagihan float64
	TotalBayar   float64
}
