package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID 			uuid.UUID 	`db:"id" json:"id"`
	Username 	string 		`db:"username" json:"username"`
	Name 		string 		`db:"name" json:"name"`
	Email 		string 		`db:"email" json:"email"`
	Password 	string 		`db:"password" json:"password"`
	Role 		string 		`db:"role" json:"role"`
	Created_At 	time.Time 	`db:"created_at" json:"created_at"`
	Updated_At 	time.Time 	`db:"updated_at" json:"updated_at"`
}