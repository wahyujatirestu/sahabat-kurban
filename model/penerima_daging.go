package model

import (
	"time"

	"github.com/google/uuid"
)

type PenerimaDaging struct {
	ID        		uuid.UUID		`db:"id"`
	Name      		string			`db:"nama"`
	Alamat    		*string			`db:"alamat"`
	Phone     		*string			`db:"phone"`
	Status    		string			`db:"status"`
	PekurbanID 		*uuid.UUID		`db:"pekurban_id"`
	Created_At 		time.Time		`db:"created_at"`
	Updated_At 		time.Time		`db:"updated_at"`
}
