package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type DistribusiDagingRepository interface {
	Create(ctx context.Context, d *model.DistribusiDaging) error
	GetAll(ctx context.Context) ([]*model.DistribusiDaging, error)	
	GetByID(ctx context.Context, id uuid.UUID) (*model.DistribusiDaging, error)
	FindByPenerimaID(ctx context.Context, penerimaID uuid.UUID) (*model.DistribusiDaging, error)
	CountTotalPaket(ctx context.Context) (int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type distribusiDagingRepository struct {
	db *sql.DB
}

func NewDistribusiDagingRepository(db *sql.DB) DistribusiDagingRepository {
	return &distribusiDagingRepository{db: db}
}

func (r *distribusiDagingRepository) Create(ctx context.Context, d *model.DistribusiDaging) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO distribusi_daging (id, penerima_id, hewan_id, jumlah_paket, tanggal_distribusi, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, d.ID, d.PenerimaID, d.HewanID, d.JumlahPaket, d.TanggalDistribusi, d.Created_At, d.Updated_At)
	return err
}

func (r *distribusiDagingRepository) GetAll(ctx context.Context) ([]*model.DistribusiDaging, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, penerima_id, hewan_id, jumlah_paket, tanggal_distribusi, created_at, updated_at FROM distribusi_daging`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.DistribusiDaging
	for rows.Next(){
		var d model.DistribusiDaging
		if err:= rows.Scan(&d.ID, &d.PenerimaID, &d.HewanID, &d.JumlahPaket, &d.TanggalDistribusi, &d.Created_At, &d.Updated_At); err != nil {
			return nil, err
		}

		result = append(result, &d)
	}

	return result, nil
}

func (r *distribusiDagingRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.DistribusiDaging, error) {
	rows := r.db.QueryRowContext(ctx, `SELECT id, penerima_id, hewan_id, jumlah_paket, tanggal_distribusi, created_at, updated_at FROM distribusi_daging WHERE id = $1`, id)

	var d model.DistribusiDaging
	if err := rows.Scan(&d.ID, &d.PenerimaID, &d.HewanID, &d.JumlahPaket, &d.TanggalDistribusi, &d.Created_At, &d.Updated_At); err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *distribusiDagingRepository) FindByPenerimaID(ctx context.Context, penerimaID uuid.UUID) (*model.DistribusiDaging, error) {
	rows := r.db.QueryRowContext(ctx, `SELECT id, penerima_id, hewan_id, jumlah_paket, tanggal_distribusi, created_at, updated_at FROM distribusi_daging WHERE penerima_id = $1`, penerimaID)

	var d model.DistribusiDaging
	if err := rows.Scan(&d.ID, &d.PenerimaID, &d.HewanID, &d.JumlahPaket, &d.TanggalDistribusi, &d.Created_At, &d.Updated_At); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

func (r *distribusiDagingRepository) CountTotalPaket(ctx context.Context) (int, error) {
	var total int
	query := `SELECT COALESCE(SUM(jumlah_paket), 0) FROM distribusi_daging`
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	return total, err
}


func (r *distribusiDagingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM distribusi_daging WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Distribusi daging not found")
	}

	return nil
}