package dto

import (
	"time"

	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type CreateDistribusiRequest struct {
	PenerimaID        string    `json:"penerima_id" binding:"required,uuid"`
	HewanID           string    `json:"hewan_id" binding:"required,uuid"`
	JumlahPaket       int       `json:"jumlah_paket" binding:"required,min=1"`
	TanggalDistribusi string `json:"tanggal_distribusi" binding:"required"`
}

type DistribusiResponse struct {
	ID                string    `json:"id"`
	PenerimaID        string    `json:"penerima_id"`
	HewanID           string    `json:"hewan_id"`
	JumlahPaket       int       `json:"jumlah_paket"`
	TanggalDistribusi time.Time `json:"tanggal_distribusi"`
}


func ToDistribusiResponse(d *model.DistribusiDaging) DistribusiResponse {
	return DistribusiResponse{
		ID: d.ID.String(),
		PenerimaID: d.PenerimaID.String(),
		HewanID: d.HewanID.String(),
		JumlahPaket: d.JumlahPaket,
		TanggalDistribusi: d.TanggalDistribusi,
	}
}