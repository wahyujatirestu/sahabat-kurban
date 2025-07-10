package model

import "github.com/google/uuid"

type PekurbanHewan struct {
	PekurbanID	uuid.UUID	`db:"pekurban_id"`
	HewanID 	uuid.UUID	`db:"hewan_id"`
	Porsi		float64		`db:"porsi"`
}