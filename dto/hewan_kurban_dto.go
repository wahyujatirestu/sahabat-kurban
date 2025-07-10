package dto

import (
	"time"

	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type CreateHewanKurbanRequest struct {
	Jenis              	string  		`json:"jenis" binding:"required,oneof=sapi kambing domba"`
	Berat              	float64 		`json:"berat" binding:"required,gt=0"`
	TglPendaftaran 		string  		`json:"tanggal_pendaftaran" binding:"required"` 
}

type UpdateHewanKurbanRequest struct {
	Jenis              	string  		`json:"jenis" binding:"omitempty,oneof=sapi kambing domba"`
	Berat              	float64 		`json:"berat" binding:"omitempty,gt=0"`
	TglPendaftaran 		string  		`json:"tanggal_pendaftaran" binding:"omitempty"` 
}

type HewanKurbanResponse struct {
	ID                 	string  		`json:"id"`
	Jenis              	string  		`json:"jenis"`
	Berat              	float64 		`json:"berat"`
	TglPendaftaran 		string  		`json:"tanggal_pendaftaran"`
	CreatedAt          	string  		`json:"created_at"`
	UpdatedAt          	string  		`json:"updated_at"`
}

func ToHewanKurbanResponse(h *model.HewanKurban) HewanKurbanResponse {
	return HewanKurbanResponse{
		ID: h.ID.String(),
		Jenis: h.Jenis,
		Berat: h.Berat,
		TglPendaftaran: h.TglPendaftaran.Format("2006-01-02"),
		CreatedAt: h.Created_At.Format(time.RFC3339),
		UpdatedAt: h.Updated_At.Format(time.RFC3339),
	}
}
