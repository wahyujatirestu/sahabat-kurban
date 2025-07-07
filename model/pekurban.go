package model

import (
	"time"

	"github.com/google/uuid"
)

type Pekurban struct {
	ID 			uuid.UUID	`db:"id" json:"id"`
	UserId		*uuid.UUID	`db:"user_id" json:"user_id,omitempty"`
	Name 		*string		`db:"name" json:"name,omitempty"`
	Phone 		*string		`db:"phone" json:"phone,omitempty"`
	Email 		*string		`db:"email" json:"email,omitempty"`
	Alamat		*string		`db:"alamat" json:"alamat,omitempty"`
	Created_At	time.Time	`db:"created_at" json:"created_at"`
	Updated_At	time.Time	`db:"updated_at" json:"updated_at"`
}