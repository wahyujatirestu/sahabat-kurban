package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/utils/model"
)

type EmailVerificationRepository interface {
	Save(ctx context.Context, token *model.EmailVerificationToken) error
	FindByToken(ctx context.Context, token string) (*model.EmailVerificationToken, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type emailVerificationRepo struct {
	db *sql.DB
}

func NewEmailVerificationRepository(db *sql.DB) EmailVerificationRepository {
	return &emailVerificationRepo{db: db}
}

func (r *emailVerificationRepo) Save(ctx context.Context, t *model.EmailVerificationToken) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO email_verification_tokens (id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`,
		t.ID, t.UserID, t.Token, t.Expires_At, t.Created_At)
	return err
}

func (r *emailVerificationRepo) FindByToken(ctx context.Context, token string) (*model.EmailVerificationToken, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token, expires_at, created_at FROM email_verification_tokens WHERE token = $1`, token)

	var t model.EmailVerificationToken
	err := row.Scan(&t.ID, &t.UserID, &t.Token, &t.Expires_At, &t.Created_At)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *emailVerificationRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM email_verification_tokens WHERE user_id = $1`, userID)
	return err
}
