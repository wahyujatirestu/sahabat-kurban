package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/utils/model"
)

type ResetPasswordRepository interface {
	Save(ctx context.Context, token *model.ResetPasswordToken) error
	FindByToken(ctx context.Context, token string) (*model.ResetPasswordToken, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type resetPasswordRepository struct {
	db *sql.DB
}

func NewResetPasswordRepository(db *sql.DB) ResetPasswordRepository {
	return &resetPasswordRepository{db}
}

func (r *resetPasswordRepository) Save(ctx context.Context, token *model.ResetPasswordToken) error {
	query := `INSERT INTO reset_password_tokens (id, user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)
	return err
}

func (r *resetPasswordRepository) FindByToken(ctx context.Context, token string) (*model.ResetPasswordToken, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM reset_password_tokens WHERE token=$1`
	row := r.db.QueryRowContext(ctx, query, token)

	var rp model.ResetPasswordToken
	if err := row.Scan(&rp.ID, &rp.UserID, &rp.Token, &rp.ExpiresAt, &rp.CreatedAt); err != nil {
		return nil, err
	}
	return &rp, nil
}

func (r *resetPasswordRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM reset_password_tokens WHERE user_id=$1", userID)
	return err
}
