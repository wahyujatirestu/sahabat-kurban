package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID 			uuid.UUID 	`db:"id"`
	Username 	string 		`db:"username"`
	Name 		string 		`db:"name"`
	Email 		string 		`db:"email"`
	Password 	string 		`db:"password"`
	Role 		string 		`db:"role"`
	Created_At 	time.Time 	`db:"created_at"`
	Updated_At 	time.Time 	`db:"updated_at"`
}