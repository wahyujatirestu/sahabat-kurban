package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID				uuid.UUID		`db:"refresh_token"`
	UserID			uuid.UUID		`db:"user_id"`
	Token			string			`db:"token"`
	Expires_At		time.Time		`db:"expires_at"`
	Revoked			bool			`db:"revoked"`
	Revoked_At		*time.Time		`db:"revoked_at,omitempty"`
	ReplaceByToken 	*string			`db:"replaced_by_token,omitempty"`
	Created_At		time.Time		`db:"created_at"`
	Updated_At		time.Time		`db:"updated_at"`
}