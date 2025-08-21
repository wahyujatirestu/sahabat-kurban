package dto

import "github.com/google/uuid"

type RekapDanaHewanResponse struct {
	HewanID     uuid.UUID `json:"hewan_id"`
	Jenis       string    `json:"jenis"`
	HargaTarget float64   `json:"harga_target"`
	TotalMasuk  float64   `json:"total_masuk"`
	Status      string    `json:"status"`
	IsPrivate   bool      `json:"is_private"`
}


type ProgressPembayaranPekurban struct {
	PekurbanID    uuid.UUID `json:"pekurban_id"`
	NamaPekurban  string    `json:"nama_pekurban"`
	JumlahPorsi   float64   `json:"jumlah_porsi"`
	TotalTagihan  float64   `json:"total_tagihan"`
	TotalBayar    float64   `json:"total_bayar"`
	Progress      float64   `json:"progress"`
	Status        string    `json:"status"`
}

