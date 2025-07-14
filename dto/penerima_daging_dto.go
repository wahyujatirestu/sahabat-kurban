package dto

import "github.com/wahyujatirestu/sahabat-kurban/model"

type CreatePenerimaRequest struct {
	Name       string  `json:"name"`
	Alamat     *string `json:"alamat,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Status     string  `json:"status" binding:"required,oneof=warga dhuafa panitia pekurban"`
	PekurbanID *string `json:"pekurban_id,omitempty"`
}

type UpdatePenerimaRequest struct {
	Name       *string `json:"name,omitempty"`
	Alamat     *string `json:"alamat,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Status     string  `json:"status" binding:"required,oneof=warga dhuafa panitia pekurban"`
	PekurbanID *string `json:"pekurban_id,omitempty"`
}

type PenerimaResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Alamat     *string `json:"alamat,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Status     string  `json:"status"`
	PekurbanID *string `json:"pekurban_id,omitempty"`
}


func ToPenerimaResponse(p *model.PenerimaDaging) PenerimaResponse {
	var pekurbanID	*string
	if p.PekurbanID != nil {
		id := p.PekurbanID.String()
		pekurbanID = &id
	}

	return 	PenerimaResponse{
		ID: p.ID.String(),
		Name: p.Name,
		Alamat: p.Alamat,
		Phone: p.Phone,
		Status: p.Status,
		PekurbanID: pekurbanID,
	}
}