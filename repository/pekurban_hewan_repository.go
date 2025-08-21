package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type PekurbanHewanRepository interface {
	Create(ctx context.Context, ph *model.PekurbanHewan) error
	FindAll(ctx context.Context) ([]*model.PekurbanHewanJoin, error)
	GetByHewanId(ctx context.Context, hewanID uuid.UUID) ([]*model.PekurbanHewanJoin, error)
	GetByPekurbanId(ctx context.Context, pekurbanID uuid.UUID) ([]*model.PekurbanHewanJoin, error)
	Update(ctx context.Context, ph *model.PekurbanHewan) error
	Delete(ctx context.Context, pekurbanID, hewanID uuid.UUID) error
}

type pekurbanHewanRepository struct {
	db *sql.DB
}

func NewPekurbanHewanRepository(db *sql.DB) PekurbanHewanRepository {
	return &pekurbanHewanRepository{db: db}
}

func (r *pekurbanHewanRepository) Create(ctx context.Context, ph *model.PekurbanHewan) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO pekurban_hewan (pekurban_id, hewan_id, porsi) VALUES ($1, $2, $3)`, ph.PekurbanID, ph.HewanID, ph.Porsi)
	return err
}

func (r *pekurbanHewanRepository) FindAll(ctx context.Context) ([]*model.PekurbanHewanJoin, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ph.pekurban_id, p.name, ph.hewan_id, h.jenis, ph.porsi
        FROM pekurban_hewan ph
        JOIN pekurban p ON ph.pekurban_id = p.id
        JOIN hewan_kurban h ON ph.hewan_id = h.id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.PekurbanHewanJoin
	for rows.Next() {
		ph, err := scanPekurbanHewan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, ph)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *pekurbanHewanRepository) GetByHewanId(ctx context.Context, hewanID uuid.UUID) ([]*model.PekurbanHewanJoin, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ph.pekurban_id, p.name, ph.hewan_id, h.jenis, ph.porsi
        FROM pekurban_hewan ph
        JOIN pekurban p ON ph.pekurban_id = p.id
        JOIN hewan_kurban h ON ph.hewan_id = h.id 
		WHERE hewan_id = $1`, hewanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.PekurbanHewanJoin
	for rows.Next() {
		ph, err := scanPekurbanHewan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, ph)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *pekurbanHewanRepository) GetByPekurbanId(ctx context.Context, pekurbanID uuid.UUID) ([]*model.PekurbanHewanJoin, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ph.pekurban_id, p.name, ph.hewan_id, h.jenis, ph.porsi
        FROM pekurban_hewan ph
        JOIN pekurban p ON ph.pekurban_id = p.id
        JOIN hewan_kurban h ON ph.hewan_id = h.id 
		WHERE pekurban_id = $1`, pekurbanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.PekurbanHewanJoin
	for rows.Next() {
		ph, err := scanPekurbanHewan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, ph)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *pekurbanHewanRepository) Update(ctx context.Context, ph *model.PekurbanHewan) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE pekurban_hewan SET porsi=$3 WHERE pekurban_id=$1 AND hewan_id=$2`,
		ph.PekurbanID, ph.HewanID, ph.Porsi)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("data not found")
	}

	return nil
}

func (r *pekurbanHewanRepository) Delete(ctx context.Context, pekurbanID, hewanID uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM pekurban_hewan WHERE pekurban_id=$1 AND hewan_id=$2`, pekurbanID, hewanID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Data not found or already deleted")
	}

	return nil
}

func scanPekurbanHewan(rows *sql.Rows) (*model.PekurbanHewanJoin, error) {
	var ph model.PekurbanHewanJoin
	err := rows.Scan(
		&ph.PekurbanID,
		&ph.Pekurban,
		&ph.HewanID,
		&ph.Hewan,
		&ph.Porsi,
	)
	if err != nil {
		return nil, err
	}
	return &ph, nil
}