package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type PenyembelihanRepository interface {
	Create(ctx context.Context, p *model.Penyembelihan) error
	GetAll(ctx context.Context) ([]*model.Penyembelihan, error)
	GetById(ctx context.Context, id uuid.UUID) (*model.Penyembelihan, error)
	GetByHewanID(ctx context.Context, hewanID uuid.UUID) (*model.Penyembelihan, error)
	Update(ctx context.Context, p *model.Penyembelihan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type penyembelihanRepository struct{
	db *sql.DB
}

func NewPenyembelihanRepository(db *sql.DB) PenyembelihanRepository {
	return &penyembelihanRepository{db: db}
}

func (r *penyembelihanRepository) Create(ctx context.Context, p *model.Penyembelihan) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO penyembelihan (id, hewan_id, tanggal_penyembelihan, lokasi, urutan_rencana, urutan_aktual, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, p.ID, p.HewanID, p.TglPenyembelihan, p.Lokasi,
		p.UrutanRencana, p.UrutanAktual, p.Created_At, p.Updated_At)
	return err
}

func (r *penyembelihanRepository) GetAll(ctx context.Context) ([]*model.Penyembelihan, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, hewan_id, tanggal_penyembelihan, lokasi, urutan_rencana, urutan_aktual, created_at, updated_at FROM penyembelihan`)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	var result []*model.Penyembelihan
	for rows.Next(){
		var p model.Penyembelihan
		if err := rows.Scan(&p.ID, &p.HewanID, &p.TglPenyembelihan, &p.Lokasi,
		&p.UrutanRencana, &p.UrutanAktual, &p.Created_At, &p.Updated_At); err != nil {
			return nil, err
		}

		result = append(result, &p)
	}

	return result, nil
}

func (r *penyembelihanRepository) GetByHewanID(ctx context.Context, hewanID uuid.UUID) (*model.Penyembelihan, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, hewan_id, tanggal_penyembelihan, lokasi, urutan_rencana, urutan_aktual, created_at, updated_at FROM penyembelihan WHERE hewan_id = $1 LIMIT 1`, hewanID)

	var p model.Penyembelihan
	if err := row.Scan(&p.ID, &p.HewanID, &p.TglPenyembelihan, &p.Lokasi, &p.UrutanRencana, &p.UrutanAktual, &p.Created_At, &p.Updated_At);err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *penyembelihanRepository) GetById(ctx context.Context, id uuid.UUID) (*model.Penyembelihan, error) {
	rows := r.db.QueryRowContext(ctx, `SELECT id, hewan_id, tanggal_penyembelihan, lokasi, urutan_rencana, urutan_aktual, created_at, updated_at FROM penyembelihan WHERE id = $1`, id)

	var p model.Penyembelihan
	if err := rows.Scan(&p.ID, &p.HewanID, &p.TglPenyembelihan, &p.Lokasi,
		&p.UrutanRencana, &p.UrutanAktual, &p.Created_At, &p.Updated_At); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	return &p, nil
}

func (r *penyembelihanRepository) Update(ctx context.Context, p *model.Penyembelihan) error {
	_, err := r.db.ExecContext(ctx, `UPDATE penyembelihan SET tanggal_penyembelihan=$2, lokasi=$3, urutan_rencana=$4, urutan_aktual=$5 WHERE id = $1`, p.ID, p.TglPenyembelihan, p.Lokasi, p.UrutanRencana, p.UrutanAktual)
	return err
}

func (r *penyembelihanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM penyembelihan WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("ID penyembelihan not found")
	}

	return nil
}