package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/wahyujatirestu/sahabat-kurban/model"
)

type ReportRepository interface {
	// blok pekurban
	GetPekurbanAggregates(ctx context.Context, f model.ReportFilter) ([]model.PekurbanAggregate, error)
	GetPekurbanHewanDetails(ctx context.Context, f model.ReportFilter) ([]model.PekurbanHewanDetail, error)

	// rekap
	GetHewanAggregate(ctx context.Context, f model.ReportFilter) ([]model.HewanAggregate, error)
	GetDistribusiAggregate(ctx context.Context, f model.ReportFilter) ([]model.DistribusiAggregate, error)
	GetPembayaranAggregate(ctx context.Context, f model.ReportFilter) ([]model.PembayaranAggregate, error)

	// summary counts
	CountPekurban(ctx context.Context) (int, error)
	CountHewan(ctx context.Context, f model.ReportFilter) (int, error)
	CountPenerima(ctx context.Context) (int, error)
	SumPaketDistribusi(ctx context.Context, f model.ReportFilter) (int, error)
	SumPembayaranSettlement(ctx context.Context, f model.ReportFilter) (float64, error)
}

type reportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return &reportRepository{db: db}
}

// ================= Helper =================

// helper untuk where clause tanggal
func betweenClause(col string, f model.ReportFilter, args *[]any) string {
	clauses := []string{}
	if f.StartDate != nil {
		clauses = append(clauses, fmt.Sprintf("%s >= $%d", col, len(*args)+1))
		*args = append(*args, f.StartDate)
	}
	if f.EndDate != nil {
		// pakai <= endDate + 1day untuk inklusif
		end := f.EndDate.Truncate(24 * time.Hour).Add(24 * time.Hour)
		clauses = append(clauses, fmt.Sprintf("%s < $%d", col, len(*args)+1))
		*args = append(*args, end)
	}
	if len(clauses) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(clauses, " AND ")
}

// helper untuk nambah kondisi tambahan dengan aman
func appendCondition(baseWhere string, condition string) string {
	if baseWhere == "" {
		return " WHERE " + condition
	}
	return baseWhere + " AND " + condition
}

// ================= Pekurban =================

func (r *reportRepository) GetPekurbanAggregates(ctx context.Context, f model.ReportFilter) ([]model.PekurbanAggregate, error) {
	args := []any{}
	pembayaranWhere := betweenClause("pk.transaction_time", f, &args)
	if pembayaranWhere == "" {
		pembayaranWhere = betweenClause("pk.tanggal_pembayaran", f, &args)
	}
	// filter hanya settlement
	pembayaranWhere = appendCondition(pembayaranWhere, "pk.status = 'settlement'")

	args2 := []any{}
	hewanWhere := betweenClause("hk.tanggal_pendaftaran", f, &args2)

	query := fmt.Sprintf(`
	WITH payment AS (
		SELECT pk.pekurban_id, COALESCE(SUM(pk.jumlah),0) AS total_bayar
		FROM pembayaran_kurban pk
		%s
		GROUP BY pk.pekurban_id
	), kewajiban AS (
		SELECT ph.pekurban_id,
		       COALESCE(SUM( CASE WHEN hk.is_private = FALSE THEN (ph.porsi * hk.harga) ELSE 0 END ),0) AS total_kewajiban,
		       COUNT(DISTINCT ph.hewan_id) AS total_hewan
		FROM pekurban_hewan ph
		JOIN hewan_kurban hk ON hk.id = ph.hewan_id
		%s
		GROUP BY ph.pekurban_id
	)
	SELECT p.id, p.name, COALESCE(p.phone,''), COALESCE(p.email,''),
	       COALESCE(k.total_hewan,0),
	       COALESCE(pay.total_bayar,0),
	       COALESCE(k.total_kewajiban,0)
	FROM pekurban p
	LEFT JOIN kewajiban k ON k.pekurban_id = p.id
	LEFT JOIN payment  pay ON pay.pekurban_id = p.id
	ORDER BY p.created_at ASC
	`, pembayaranWhere, hewanWhere)

	// gabungkan args
	args = append(args, args2...)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.PekurbanAggregate{}
	for rows.Next() {
		var a model.PekurbanAggregate
		if err := rows.Scan(&a.PekurbanID, &a.Name, &a.Phone, &a.Email, &a.TotalHewan, &a.TotalBayar, &a.TotalKewajiban); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *reportRepository) GetPekurbanHewanDetails(ctx context.Context, f model.ReportFilter) ([]model.PekurbanHewanDetail, error) {
	args := []any{}
	where := betweenClause("hk.tanggal_pendaftaran", f, &args)

	q := fmt.Sprintf(`
	SELECT ph.pekurban_id, hk.id, hk.jenis::text, COALESCE(hk.berat,0), COALESCE(hk.harga,0), ph.porsi
	FROM pekurban_hewan ph
	JOIN hewan_kurban hk ON hk.id = ph.hewan_id
	%s
	ORDER BY ph.pekurban_id
	`, where)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.PekurbanHewanDetail{}
	for rows.Next() {
		var d model.PekurbanHewanDetail
		if err := rows.Scan(&d.PekurbanID, &d.HewanID, &d.Jenis, &d.Berat, &d.Harga, &d.Porsi); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// ================= Rekap =================

func (r *reportRepository) GetHewanAggregate(ctx context.Context, f model.ReportFilter) ([]model.HewanAggregate, error) {
	args := []any{}
	where := betweenClause("hk.tanggal_pendaftaran", f, &args)

	q := fmt.Sprintf(`
	SELECT hk.jenis::text, COUNT(*) AS total, COALESCE(SUM(hk.berat),0) AS berat, COALESCE(SUM(hk.harga),0) AS harga
	FROM hewan_kurban hk
	%s
	GROUP BY hk.jenis
	ORDER BY hk.jenis
	`, where)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.HewanAggregate{}
	for rows.Next() {
		var h model.HewanAggregate
		if err := rows.Scan(&h.Jenis, &h.Total, &h.Berat, &h.Harga); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}

func (r *reportRepository) GetDistribusiAggregate(ctx context.Context, f model.ReportFilter) ([]model.DistribusiAggregate, error) {
	args := []any{}
	where := betweenClause("dd.tanggal_distribusi", f, &args)

	q := fmt.Sprintf(`
	SELECT pd.status::text, COUNT(DISTINCT pd.id) AS total_penerima,
	       COALESCE(SUM(dd.jumlah_paket),0) AS total_paket
	FROM penerima_daging pd
	LEFT JOIN distribusi_daging dd ON dd.penerima_id = pd.id
	%s
	GROUP BY pd.status
	ORDER BY pd.status
	`, where)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.DistribusiAggregate{}
	for rows.Next() {
		var d model.DistribusiAggregate
		if err := rows.Scan(&d.Status, &d.TotalPenerima, &d.TotalPaket); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *reportRepository) GetPembayaranAggregate(ctx context.Context, f model.ReportFilter) ([]model.PembayaranAggregate, error) {
	args := []any{}
	where := betweenClause("pk.transaction_time", f, &args)
	if where == "" {
		where = betweenClause("pk.tanggal_pembayaran", f, &args)
	}

	q := fmt.Sprintf(`
	SELECT pk.status, COUNT(*) AS total, COALESCE(SUM(pk.jumlah),0) AS jumlah
	FROM pembayaran_kurban pk
	%s
	GROUP BY pk.status
	ORDER BY pk.status
	`, where)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.PembayaranAggregate{}
	for rows.Next() {
		var p model.PembayaranAggregate
		if err := rows.Scan(&p.Status, &p.Total, &p.Jumlah); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ================= Summary counts =================

func (r *reportRepository) CountPekurban(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pekurban`).Scan(&n)
	return n, err
}

func (r *reportRepository) CountHewan(ctx context.Context, f model.ReportFilter) (int, error) {
	args := []any{}
	where := betweenClause("tanggal_pendaftaran", f, &args)
	q := fmt.Sprintf(`SELECT COUNT(*) FROM hewan_kurban %s`, where)
	var n int
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&n)
	return n, err
}

func (r *reportRepository) CountPenerima(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM penerima_daging`).Scan(&n)
	return n, err
}

func (r *reportRepository) SumPaketDistribusi(ctx context.Context, f model.ReportFilter) (int, error) {
	args := []any{}
	where := betweenClause("tanggal_distribusi", f, &args)
	q := fmt.Sprintf(`SELECT COALESCE(SUM(jumlah_paket),0) FROM distribusi_daging %s`, where)
	var n int
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&n)
	return n, err
}

func (r *reportRepository) SumPembayaranSettlement(ctx context.Context, f model.ReportFilter) (float64, error) {
	args := []any{}
	where := betweenClause("transaction_time", f, &args)
	if where == "" {
		where = betweenClause("tanggal_pembayaran", f, &args)
	}
	where = appendCondition(where, "pk.status='settlement'")

	q := fmt.Sprintf(`SELECT COALESCE(SUM(jumlah),0) FROM pembayaran_kurban pk %s`, where)

	var sum float64
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&sum)
	return sum, err
}
