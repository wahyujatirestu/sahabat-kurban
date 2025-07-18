package dto

import (
	"time"

	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type CreateHewanKurbanRequest struct {
	Jenis           string  `json:"jenis" binding:"required,oneof=sapi kambing domba"`
	Berat           float64 `json:"berat" binding:"required,gt=0"`
	Harga           *float64 `json:"harga"`
	IsPrivate       *bool   `json:"is_private"`
	TglPendaftaran  string  `json:"tanggal_pendaftaran" binding:"required"`
}

type UpdateHewanKurbanRequest struct {
	Jenis           string   `json:"jenis" binding:"omitempty,oneof=sapi kambing domba"`
	Berat           float64  `json:"berat" binding:"omitempty,gt=0"`
	Harga           float64  `json:"harga" binding:"omitempty,gt=0"`
	IsPrivate       *bool    `json:"is_private"`
	TglPendaftaran  string   `json:"tanggal_pendaftaran" binding:"omitempty"`
}

type HewanKurbanResponse struct {
	ID              	string  `json:"id"`
	Jenis           	string  `json:"jenis"`
	Berat           	float64 `json:"berat"`
	Harga           	float64 `json:"harga"`
	IsPrivate       	bool    `json:"is_private"`
	TglPendaftaran  	string  `json:"tanggal_pendaftaran"`
	StatusPenyembelihan string  `json:"status_penyembelihan"`
	CreatedAt       	string  `json:"created_at"`
	UpdatedAt       	string  `json:"updated_at"`
}

func ToHewanKurbanResponse(h *model.HewanKurban, sudahDisembelih bool) HewanKurbanResponse {
	status := "belum"
	if sudahDisembelih {
		status = "sudah"
	}

	return HewanKurbanResponse{
		ID:             h.ID.String(),
		Jenis:          string(h.Jenis),
		Berat:          h.Berat,
		Harga:          h.Harga,
		IsPrivate:      h.IsPrivate,
		TglPendaftaran: h.TanggalPendaftaran.Format("2006-01-02"),
		StatusPenyembelihan: status,
		CreatedAt:      h.Created_At.Format(time.RFC3339),
		UpdatedAt:      h.Updated_At.Format(time.RFC3339),
	}
}
