package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type PenerimaDagingRepository interface {
	Create(ctx context.Context, p *model.PenerimaDaging) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.PenerimaDaging, error)
	GetAll(ctx context.Context) ([]*model.PenerimaDaging, error)
	Update(ctx context.Context, p *model.PenerimaDaging) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type penerimaDagingRepository struct {
	db *sql.DB
}

func NewPenerimaDagingRepository(db *sql.DB) PenerimaDagingRepository {
	return &penerimaDagingRepository{db}
}

func (r *penerimaDagingRepository) Create(ctx context.Context, p *model.PenerimaDaging ) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO penerima_daging (id, name, alamat, phone, status, pekurban_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, p.ID, p.Name, p.Alamat, p.Phone, p.Status, p.PekurbanID, p.Created_At, p.Updated_At)
	return err
}

func (r *penerimaDagingRepository) Update(ctx context.Context, p *model.PenerimaDaging) error {
	_, err := r.db.ExecContext(ctx, `UPDATE penerima_daging SET name=$2, alamat=$3, phone=$4, status=$5, pekurban_id=$6 WHERE id = $1`, p.ID, p.Name, p.Alamat, p.Phone, p.Status, p.PekurbanID)
	return err
}

func (r *penerimaDagingRepository) GetAll(ctx context.Context) ([]*model.PenerimaDaging, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, alamat, phone, status, pekurban_id, created_at, updated_at FROM penerima_daging`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.PenerimaDaging
	for rows.Next() {
		var p model.PenerimaDaging
		if err := rows.Scan(&p.ID, &p.Name, &p.Alamat, &p.Phone, &p.Status, &p.PekurbanID, &p.Created_At, &p.Updated_At); err != nil {
			return nil, err
		}

		result = append(result, &p)
	}

	return result, nil
}

func (r *penerimaDagingRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PenerimaDaging, error) {
	rows := r.db.QueryRowContext(ctx, "SELECT id, name, alamat, phone, status, pekurban_id, created_at, updated_at FROM penerima_daging WHERE id = $1", id)

	var p model.PenerimaDaging
	if err := rows.Scan(&p.ID, &p.Name, &p.Alamat, &p.Phone, &p.Status, &p.PekurbanID, &p.Created_At, &p.Updated_At); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (r *penerimaDagingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM penerima_daging WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Penerima daging not found")
	}

	return nil
}