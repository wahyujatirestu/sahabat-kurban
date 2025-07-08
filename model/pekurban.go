package model

import (
	"time"

	"github.com/google/uuid"
)

type Pekurban struct {
	ID 			uuid.UUID	`db:"id"`
	UserId		*uuid.UUID	`db:"user_id"`
	Name 		*string		`db:"name"`
	Phone 		*string		`db:"phone"`
	Email 		*string		`db:"email"`
	Alamat		*string		`db:"alamat"`
	Created_At	time.Time	`db:"created_at"`
	Updated_At	time.Time	`db:"updated_at"`
}