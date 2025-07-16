// repository/pembayaran_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type PembayaranKurbanRepository interface {
	Create(ctx context.Context, p *model.PembayaranKurban) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.PembayaranKurban, error)
	FindByOrderID(ctx context.Context, orderID string) (*model.PembayaranKurban, error)
	GetAll(ctx context.Context) ([]*model.PembayaranKurban, error)
}

type pembayaranRepo struct {
	db *sql.DB
}

func NewPembayaranKurbanRepository(db *sql.DB) PembayaranKurbanRepository {
	return &pembayaranRepo{db: db}
}

func (r *pembayaranRepo) Create(ctx context.Context, p *model.PembayaranKurban) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO pembayaran_kurban (
		id, order_id, transaction_id, pekurban_id, metode, payment_type, va_number, jumlah, status, fraud_status,
		approval_code, transaction_time, tanggal_pembayaran, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		p.ID, p.OrderID, p.TransactionID, p.PekurbanID, p.Metode, p.PaymentType, p.VANumber,
		p.Jumlah, p.Status, p.FraudStatus, p.ApprovalCode, p.TransactionTime, p.TanggalPembayaran,
		p.Created_At, p.Updated_At,
	)
	return err
}

func (r *pembayaranRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.PembayaranKurban, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, order_id, transaction_id, pekurban_id, metode, payment_type, va_number,
		jumlah, status, fraud_status, approval_code, transaction_time, tanggal_pembayaran, created_at, updated_at
		FROM pembayaran_kurban WHERE id = $1`, id)
	return scanPembayaran(row)
}

func (r *pembayaranRepo) FindByOrderID(ctx context.Context, orderID string) (*model.PembayaranKurban, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, order_id, transaction_id, pekurban_id, metode, payment_type, va_number,
		jumlah, status, fraud_status, approval_code, transaction_time, tanggal_pembayaran, created_at, updated_at
		FROM pembayaran_kurban WHERE order_id = $1`, orderID)
	return scanPembayaran(row)
}

func (r *pembayaranRepo) GetAll(ctx context.Context) ([]*model.PembayaranKurban, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, order_id, transaction_id, pekurban_id, metode, payment_type, va_number,
		jumlah, status, fraud_status, approval_code, transaction_time, tanggal_pembayaran, created_at, updated_at
		FROM pembayaran_kurban`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.PembayaranKurban
	for rows.Next() {
		p := new(model.PembayaranKurban)
		err := rows.Scan(
			&p.ID, &p.OrderID, &p.TransactionID, &p.PekurbanID, &p.Metode, &p.PaymentType, &p.VANumber,
			&p.Jumlah, &p.Status, &p.FraudStatus, &p.ApprovalCode, &p.TransactionTime, &p.TanggalPembayaran,
			&p.Created_At, &p.Updated_At,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func scanPembayaran(row *sql.Row) (*model.PembayaranKurban, error) {
	var p model.PembayaranKurban
	err := row.Scan(
		&p.ID, &p.OrderID, &p.TransactionID, &p.PekurbanID, &p.Metode, &p.PaymentType, &p.VANumber,
		&p.Jumlah, &p.Status, &p.FraudStatus, &p.ApprovalCode, &p.TransactionTime, &p.TanggalPembayaran,
		&p.Created_At, &p.Updated_At,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
