package dto

import "github.com/wahyujatirestu/sahabat-kurban/model"

type CreatePekurbanHewanRequest struct {
	PekurbanID	 string	`json:"pekurban_id" binding:"required,uuid"`
	HewanID		 string	`json:"hewan_id" binding:"required,uuid"`
	JumlahOrang  int    `json:"jumlah_orang" binding:"required,gt=0,lte=7"`
}

type PekurbanHewanResponse struct {
	PekurbanID string  `json:"pekurban_id"`
	HewanID    string  `json:"hewan_id"`
	Porsi      float64 `json:"porsi"`
}

type UpdatePekurbanHewanRequest struct {
	JumlahOrang int `json:"jumlah_orang" binding:"required,gt=0,lte=7"`
}


func ToPekurbanHewanResponse(ph *model.PekurbanHewan) PekurbanHewanResponse {
	return PekurbanHewanResponse{
		PekurbanID: ph.PekurbanID.String(),
		HewanID: ph.HewanID.String(),
		Porsi: ph.Porsi,
	}
}