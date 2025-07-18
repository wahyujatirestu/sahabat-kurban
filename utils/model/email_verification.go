package model

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerificationToken struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	Token      string    `db:"token"`
	Expires_At  time.Time `db:"expires_at"`
	Created_At  time.Time `db:"created_at"`
}
