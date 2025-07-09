package dto

type CreatePekurbanRequest struct {
	UserID *string `json:"userId"`
	Name   string  `json:"name" binding:"required"`
	Phone  string  `json:"phone" binding:"required"`
	Email  string  `json:"email" binding:"required,email"`
	Alamat string  `json:"alamat" binding:"required"`
}

type UpdatePekurbanRequest struct {
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	Alamat string `json:"alamat" binding:"required"`
}

type PekurbanResponse struct {
	ID     string  `json:"id"`
	UserID *string `json:"user_id,omitempty"`
	Name   string  `json:"name"`
	Phone  string  `json:"phone"`
	Email  string  `json:"email"`
	Alamat string  `json:"alamat"`
}
