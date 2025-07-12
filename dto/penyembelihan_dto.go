package dto

import (
	"time"

	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type CreatePenyembelihanRequest struct {
	HewanID              string    `json:"hewan_id" binding:"required,uuid"`
	TanggalPenyembelihan time.Time `json:"tanggal_penyembelihan" binding:"required"`
	Lokasi               string    `json:"lokasi" binding:"required"`
	UrutanRencana        int       `json:"urutan_rencana" binding:"omitempty,min=1"`
}

type UpdatePenyembelihanRequest struct {
	TanggalPenyembelihan time.Time `json:"tanggal_penyembelihan" binding:"required"`
	Lokasi               string    `json:"lokasi" binding:"required"`
	UrutanRencana        int       `json:"urutan_rencana" binding:"omitempty,min=1"`
	UrutanAktual         *int      `json:"urutan_aktual"`
}

type PenyembelihanResponse struct {
	ID                   string    `json:"id"`
	HewanID              string    `json:"hewan_id"`
	TanggalPenyembelihan time.Time `json:"tanggal_penyembelihan"`
	Lokasi               string    `json:"lokasi"`
	UrutanRencana        int       `json:"urutan_rencana"`
	UrutanAktual         *int      `json:"urutan_aktual"`
}


func ToPenyembelihanResponse(p *model.Penyembelihan) PenyembelihanResponse {
	return PenyembelihanResponse{
		ID:                   p.ID.String(),
		HewanID:              p.HewanID.String(),
		TanggalPenyembelihan: p.TglPenyembelihan,
		Lokasi:               p.Lokasi,
		UrutanRencana:        p.UrutanRencana,
		UrutanAktual:         p.UrutanAktual,
	}
}
