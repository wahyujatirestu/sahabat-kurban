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
	GetTotalPembayaranPerHewan(ctx context.Context) ([]model.TotalPembayaranPerHewan, error)
	GetProgressPembayaranPekurban(ctx context.Context) ([]model.ProgressPembayaran, error)
}

type pembayaranRepo struct {
	db *sql.DB
}

func NewPembayaranKurbanRepository(db *sql.DB) PembayaranKurbanRepository {
	return &pembayaranRepo{db: db}
}

func (r *pembayaranRepo) Create(ctx context.Context, p *model.PembayaranKurban) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO pembayaran_kurban (
		id, order_id, transaction_id, pekurban_id, metode, payment_type, va_number, status, fraud_status,
		approval_code, transaction_time, tanggal_pembayaran, jumlah, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		p.ID, p.OrderID, p.TransactionID, p.PekurbanID, p.Metode, p.PaymentType, p.VANumber,
		p.Status, p.FraudStatus, p.ApprovalCode, p.TransactionTime, p.TanggalPembayaran, p.Jumlah,
		p.Created_At, p.Updated_At,
	)
	return err
}

func (r *pembayaranRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.PembayaranKurban, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, order_id, transaction_id, pekurban_id, metode, payment_type, va_number,
		status, fraud_status, approval_code, transaction_time, tanggal_pembayaran, jumlah, created_at, updated_at
		FROM pembayaran_kurban WHERE id = $1`, id)
	return scanPembayaran(row)
}

func (r *pembayaranRepo) FindByOrderID(ctx context.Context, orderID string) (*model.PembayaranKurban, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, order_id, transaction_id, pekurban_id, metode, payment_type, va_number,
		status, fraud_status, approval_code, transaction_time, tanggal_pembayaran, jumlah, created_at, updated_at
		FROM pembayaran_kurban WHERE order_id = $1`, orderID)
	return scanPembayaran(row)
}

func (r *pembayaranRepo) GetAll(ctx context.Context) ([]*model.PembayaranKurban, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, order_id, transaction_id, pekurban_id, metode, payment_type, va_number,
		status, fraud_status, approval_code, transaction_time, tanggal_pembayaran, jumlah, created_at, updated_at
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
			&p.Status, &p.FraudStatus, &p.ApprovalCode, &p.TransactionTime, &p.TanggalPembayaran, &p.Jumlah,
			&p.Created_At, &p.Updated_At,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *pembayaranRepo) GetTotalPembayaranPerHewan(ctx context.Context) ([]model.TotalPembayaranPerHewan, error) {
	query := `
	SELECT 
		h.id AS hewan_id,
		h.jenis,
		h.harga AS harga_target,
		COALESCE(SUM(ph.porsi * h.harga), 0) AS total_masuk
	FROM hewan_kurban h
	LEFT JOIN pekurban_hewan ph ON h.id = ph.hewan_id
	LEFT JOIN pembayaran_kurban pk ON ph.pekurban_id = pk.pekurban_id
		AND pk.status IN ('settlement', 'capture')
	GROUP BY h.id
	ORDER BY h.tanggal_pendaftaran;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.TotalPembayaranPerHewan
	for rows.Next() {
		var p model.TotalPembayaranPerHewan
		if err := rows.Scan(&p.HewanID, &p.Jenis, &p.HargaTarget, &p.TotalMasuk); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func (r *pembayaranRepo) GetProgressPembayaranPekurban(ctx context.Context) ([]model.ProgressPembayaran, error) {
	query := `
	SELECT 
		p.id AS pekurban_id,
		COALESCE(p.name, 'Tanpa Nama') AS nama_pekurban,
		COALESCE(SUM(ph.porsi), 0) AS total_porsi,
		COALESCE(SUM(ph.porsi * h.harga), 0) AS total_tagihan,
		COALESCE((
			SELECT SUM(h.harga * ph2.porsi)
			FROM pembayaran_kurban pk2
			JOIN pekurban_hewan ph2 ON ph2.pekurban_id = pk2.pekurban_id
			JOIN hewan_kurban h ON h.id = ph2.hewan_id
			WHERE pk2.pekurban_id = p.id
			AND pk2.status IN ('settlement', 'capture')
		), 0) AS total_bayar
	FROM pekurban p
	LEFT JOIN pekurban_hewan ph ON p.id = ph.pekurban_id
	LEFT JOIN hewan_kurban h ON h.id = ph.hewan_id
	GROUP BY p.id, p.name
	ORDER BY p.created_at
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.ProgressPembayaran
	for rows.Next() {
		var p model.ProgressPembayaran
		if err := rows.Scan(&p.PekurbanID, &p.NamaPekurban, &p.PorsiTotal, &p.TotalTagihan, &p.TotalBayar); err != nil {
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
		&p.Status, &p.FraudStatus, &p.ApprovalCode, &p.TransactionTime, &p.TanggalPembayaran, &p.Jumlah,
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

