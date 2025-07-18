package dto

import "github.com/wahyujatirestu/sahabat-kurban/model"

type CreatePekurbanRequest struct {
	UserID *string `json:"user_id"`
	Name   string  `json:"name"`
	Phone  string  `json:"phone" binding:"required"`
	Email  string  `json:"email"`
	Alamat string  `json:"alamat" binding:"required"`
}

type UpdatePekurbanRequest struct {
	Name   *string `json:"name"`
	Phone  *string `json:"phone"`
	Email  *string `json:"email"`
	Alamat *string `json:"alamat"`
}


type PekurbanResponse struct {
	ID     string  `json:"id"`
	UserID *string `json:"user_id,omitempty"`
	Name   string  `json:"name"`
	Phone  string  `json:"phone"`
	Email  string  `json:"email"`
	Alamat string  `json:"alamat"`
}

func ToPekurbanRespon(p *model.Pekurban) PekurbanResponse {
	var userID *string
	if p.UserId != nil {
		uidStr := p.UserId.String()
		userID = &uidStr
	}

	return PekurbanResponse{
		ID: p.ID.String(),
		UserID: userID,
		Name: *p.Name,
		Phone: *p.Phone,
		Email: *p.Email,
		Alamat: *p.Alamat,
	}
}
