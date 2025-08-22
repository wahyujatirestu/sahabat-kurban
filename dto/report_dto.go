package dto

import "time"

// Query param parsing di controller
type ReportQuery struct {
	StartDate *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"end_date" time_format:"2006-01-02"`
}

// Response JSON akhir
type LaporanResponse struct {
	Tanggal   		time.Time        `json:"tanggal"`
	Ringkasan 		RingkasanDTO     `json:"ringkasan"`
	Pekurban  		[]PekurbanDTO    `json:"detail_pekurban"`
	RekapHewan      []HewanDTO      `json:"rekap_hewan"`
	RekapDistribusi []DistribusiDTO `json:"distribusi"`
	RekapPembayaran []PembayaranDTO `json:"rekap_pembayaran"`
}

type RingkasanDTO struct {
	TotalPekurban             int     `json:"total_pekurban"`
	TotalHewan                int     `json:"total_hewan"`
	TotalPenerima             int     `json:"total_penerima"`
	TotalPaketDistribusi      int     `json:"total_paket_distribusi"`
	TotalPembayaranSettlement float64 `json:"total_pembayaran"`
}

type PekurbanDTO struct {
	PekurbanID       string            `json:"pekurban_id"`
	Nama             string            `json:"nama"`
	Phone            string            `json:"phone"`
	Email            string            `json:"email"`
	TotalHewan       int               `json:"total_hewan"`
	TotalBayar       float64           `json:"total_bayar"`
	TotalKewajiban   float64           `json:"total_kewajiban"`
	StatusPembayaran string            `json:"status_pembayaran"` // lunas/belum
	Hewan            []PekurbanHewanDTO `json:"hewan"`
}

type PekurbanHewanDTO struct {
	HewanID string  `json:"hewan_id"`
	Jenis   string  `json:"jenis"`
	Berat   float64 `json:"berat"`
	Harga   float64 `json:"harga"`
	Porsi   float64 `json:"porsi"`
}

type HewanDTO struct {
	Jenis string  `json:"jenis"`
	Total int     `json:"total"`
	Berat float64 `json:"berat"`
	Harga float64 `json:"harga"`
}

type DistribusiDTO struct {
	Kategori      string `json:"kategori"`
	TotalPenerima int    `json:"total_penerima"`
	TotalPaket    int    `json:"total_paket"`
}

type PembayaranDTO struct {
	Status string  `json:"status"`
	Total  int     `json:"total"`
	Jumlah float64 `json:"jumlah"`
}
