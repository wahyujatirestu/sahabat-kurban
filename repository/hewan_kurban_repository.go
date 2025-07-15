package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type HewanKurbanRepository interface {
	Create(ctx context.Context, h *model.HewanKurban) error
	GetAll(ctx context.Context) ([]*model.HewanKurban, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.HewanKurban, error)
	Update(ctx context.Context, h *model.HewanKurban) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type hewanKurbanRepository struct {
	db *sql.DB
}

func NewHewanKurbanRepository(db *sql.DB) HewanKurbanRepository {
	return &hewanKurbanRepository{db: db}
}

func (r *hewanKurbanRepository) Create(ctx context.Context, h *model.HewanKurban) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO hewan_kurban (id, jenis, berat, tanggal_pendaftaran, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`, h.ID, h.Jenis, h.Berat, h.TglPendaftaran, h.Created_At, h.Updated_At)
	return  err
}

func (r *hewanKurbanRepository) GetAll(ctx context.Context) ([]*model.HewanKurban, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, jenis, berat, tanggal_pendaftaran, created_at, updated_at FROM hewan_kurban`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*model.HewanKurban
	for rows.Next() {
		var r model.HewanKurban
		err := rows.Scan(&r.ID, &r.Jenis, &r.Berat, &r.TglPendaftaran, &r.Created_At, &r.Updated_At)
		if err != nil {
			return nil, err
		}

		result = append(result, &r)
	}

	return result, nil
}

func (r *hewanKurbanRepository) GetById(ctx context.Context, id uuid.UUID) (*model.HewanKurban, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, jenis, berat, tanggal_pendaftaran, created_at, updated_at FROM hewan_kurban WHERE id=$1`, id)

	var h model.HewanKurban
	err := row.Scan(&h.ID, &h.Jenis, &h.Berat, &h.TglPendaftaran, &h.Created_At, &h.Updated_At)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &h, nil
}

func (r *hewanKurbanRepository) Update(ctx context.Context, h *model.HewanKurban) error {
	_, err := r.db.ExecContext(ctx, `UPDATE hewan_kurban SET jenis=$2, berat=$3, tanggal_pendaftaran=$4 WHERE id=$1`, h.ID, h.Jenis, h.Berat, h.TglPendaftaran)
	return err
}

func (r *hewanKurbanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM hewan_kurban WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Hewan kurban not found")
	}

	return nil
}