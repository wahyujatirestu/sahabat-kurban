package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type PekurbanRepository interface {
	Create(ctx context.Context, p *model.Pekurban) error
	FindAll(ctx context.Context) ([]*model.Pekurban, error)
	FindById(ctx context.Context, id uuid.UUID) (*model.Pekurban, error)
	FindByUserId(ctx context.Context, userID uuid.UUID) (*model.Pekurban, error)
	Update(ctx context.Context, p *model.Pekurban) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type pekurbanRepository struct {
	db *sql.DB
}

func NewPekurbanRepository(db *sql.DB) PekurbanRepository {
	return &pekurbanRepository{db: db}
}

func (r *pekurbanRepository) Create(ctx context.Context, p *model.Pekurban) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO pekurban (id, user_id, name, phone, email, alamat, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, p.ID, p.UserId, p.Name, p.Phone, p.Email, p.Alamat, p.Created_At, p.Updated_At)
	return err
}

func (r *pekurbanRepository) FindAll(ctx context.Context) ([]*model.Pekurban, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, name, phone, email, alamat, created_at, updated_at FROM pekurban`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*model.Pekurban
	for rows.Next(){
		var r model.Pekurban
		err := rows.Scan(&r.ID, &r.UserId, &r.Name, &r.Phone, &r.Email, &r.Alamat, &r.Created_At, &r.Updated_At)
		if err != nil {
			return nil, err
		}

		result = append(result, &r)
	}

	return result, nil
}

func (r *pekurbanRepository) FindById(ctx context.Context, id uuid.UUID) (*model.Pekurban, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, name, phone, email, alamat, created_at, updated_at FROM pekurban WHERE id=$1`, id)

	var p model.Pekurban
	err := row.Scan(&p.ID, &p.UserId, &p.Name, &p.Phone, &p.Email, &p.Alamat, &p.Created_At, &p.Updated_At)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *pekurbanRepository) FindByUserId(ctx context.Context, userID uuid.UUID) (*model.Pekurban, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, name, phone, email, alamat, created_at, updated_at FROM pekurban WHERE user_id=$1`, userID)

	var p model.Pekurban
	err := row.Scan(&p.ID, &p.UserId, &p.Name, &p.Phone, &p.Email, &p.Alamat, &p.Created_At, &p.Updated_At)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *pekurbanRepository) Update(ctx context.Context, p *model.Pekurban) error {
	_, err := r.db.ExecContext(ctx, `UPDATE pekurban SET name=$1, phone=$2, email=$3, alamat=$4, updated_at=now() WHERE id=$5`, p.Name, p.Phone, p.Email, p.Alamat, p.ID)
	return err
}

func (r *pekurbanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM pekurban WHERE id=$1`, id)
	return err
}