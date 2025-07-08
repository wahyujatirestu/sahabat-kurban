package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/utils/model"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *model.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	RevokeById(ctx context.Context, id uuid.UUID, revokedAt time.Time) error
	RevokeByToken(ctx context.Context, token string, revokedAt time.Time) error
	DeleteByToken(ctx context.Context, token string) error
	DeleteById(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type refreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Save(ctx context.Context, token *model.RefreshToken) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO refresh_tokens (id, user_id, token, created_at, expires_at, revoked, revoked_at, replaced_by_token, updated_at) VALUES ($1, $2, $3, $4, $5, 6, $7, $8, $9)`, token.ID, token.UserID, token.Token, token.Created_At, token.Expires_At, token.Revoked, token.Revoked_At, token.ReplaceByToken, token.Updated_At)

	return err
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var t model.RefreshToken
	if err := r.db.QueryRowContext(ctx, `SELECT id, user_id, token, created_at, expires_at, revoked, revoked_at, replaced_by_token, updated_at FROM refresh_tokens WHERE token=$1`, token).Scan(&t.ID,
		&t.UserID,
		&t.Token,
		&t.Created_At,
		&t.Expires_At,
		&t.Revoked,
		&t.Revoked_At,
		&t.ReplaceByToken,
		&t.Updated_At,); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}

		return &t, nil
}

func (r *refreshTokenRepository) RevokeById(ctx context.Context, id uuid.UUID, revokedAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked = TRUE, revoked_at = $2, updated_at = $3, WHERE id = $1`, id, revokedAt)
	return err
}

func (r *refreshTokenRepository) RevokeByToken(ctx context.Context, token string, revokedAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked = TRUE, revoked_at = $2, updated_at = $3, WHERE token = $1`, token, revokedAt)
	return err
}

func (r *refreshTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE token = $1`, token)
	return err
}

func (r *refreshTokenRepository) DeleteById(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE expires_at < now()`)
	return err
}